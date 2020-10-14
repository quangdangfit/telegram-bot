package services

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/jinzhu/copier"
	"transport/lib/errors"
	"transport/lib/utils/logger"
	"transport/lib/utils/paging"

	"telegram-bot/app/external"
	"telegram-bot/app/models"
	"telegram-bot/app/queue"
	"telegram-bot/app/repositories"
	"telegram-bot/app/schema"
	"telegram-bot/config"
	"telegram-bot/pkg/utils"
)

const (
	MsgChanSize = 100
)

type InService interface {
	Consume()
	CronRetry(ctx context.Context) error
	CronRetryPrevious(ctx context.Context) error
	CronArchivedMessages(ctx context.Context) error
	GetMessages(ctx context.Context, query *schema.InMessageQueryParam) (*[]models.InMessage, *paging.Paging, error)
	UpdateMessage(ctx context.Context, id string, body *schema.InMessageBodyUpdateParam) (*models.InMessage, error)
}

type inService struct {
	inMsgRepo       repositories.InRepository
	routingRepo     repositories.RoutingRepository
	externalService external.External

	consumer        queue.Consumer
	consumerThreads int
}

func NewInService(externalService external.External, consumer queue.Consumer,
	inMsgRepo repositories.InRepository, routingRepo repositories.RoutingRepository) InService {
	service := inService{
		inMsgRepo:       inMsgRepo,
		routingRepo:     routingRepo,
		consumer:        consumer,
		externalService: externalService,
	}

	service.consumerThreads = config.GetConsumerThreads()

	return &service
}

func (i *inService) Consume() {
	msgChan, err := i.consumer.Consume()
	if err != nil {
		logger.Error("Cannot consume messages")
	}

	logger.Infof("Run %d threads to consume messages", i.consumerThreads)
	for index := 0; index < i.consumerThreads; index++ {
		go i.listenMessage(msgChan)
	}
}

func (i *inService) listenMessage(msgChan chan *models.InMessage) {
	for msg := range msgChan {
		i.handleMessage(context.TODO(), msg, msg.RoutingKey.Name)
		i.inMsgRepo.AddMessage(msg)
	}
}

func (i *inService) CronRetry(ctx context.Context) error {
	msgChan := make(chan models.InMessage, MsgChanSize)
	query := schema.InMessageQueryParam{
		Status: models.InMessageStatusWaitRetry,
		Page:   1,
		Limit:  config.GetRetrySize(),
		Sort:   "_id",
	}

	messages, pageInfo, _ := i.inMsgRepo.GetMessages(&query)
	if messages == nil || len(*messages) <= 0 {
		logger.Info("[Retry Message] Not found any wait_retry message!")
		return nil
	}

	logger.Infof("[Retry Message] Get %d/%d `wait_retry` messages.", len(*messages), pageInfo.Total)

	// Add `wait_retry` messages to channel
	go func() {
		for _, msg := range *messages {
			msgChan <- msg
		}
		close(msgChan)
	}()

	// handle messages from channel
	for index := 0; index < i.consumerThreads; index++ {
		go func() {
			for msg := range msgChan {
				i.retry(ctx, &msg)
			}
		}()
	}

	logger.Info("[Retry Message] Finish!")
	return nil
}

func (i *inService) retry(ctx context.Context, msg *models.InMessage) error {
	err := i.handleMessage(ctx, msg, msg.RoutingKey.Name)
	if err != nil {
		msg.Attempts += 1
		maxRetryTimes := msg.RoutingKey.RetryTimes
		if maxRetryTimes <= 0 {
			maxRetryTimes = config.DefaultMaxRetryTimes
		}
		if msg.Attempts >= maxRetryTimes {
			msg.Status = models.InMessageStatusFailed
		}
	}

	err = i.inMsgRepo.UpdateMessage(msg)
	if err != nil {
		logger.Errorf("Handle successfully, failed to update status: %s, %s, %s, error: %s", msg.RoutingKey.Name, msg.OriginModel, msg.OriginCode, err)
	}

	return err
}

func (i *inService) CronRetryPrevious(ctx context.Context) error {
	ignore := map[string]bool{}
	query := schema.InMessageQueryParam{
		Status: models.InMessageStatusWaitPrevMsg,
		Page:   1,
		Limit:  config.CronRetrySize,
		Sort:   "publish_time",
	}
	messages, _, _ := i.inMsgRepo.GetMessages(&query)
	if messages == nil {
		logger.Info("[Retry Prev Message] Not found any wait_prev message!")
		return nil
	}

	logger.Infof("[Retry Prev Message] Found %d wait_prev messages!", len(*messages))
	for _, msg := range *messages {
		key := fmt.Sprintf("%s_%s", msg.OriginModel, msg.OriginCode)
		if ignore[key] {
			logger.Info("[Retry Prev Message] Ignore message ", msg.ID)
			continue
		}

		err := i.handleMessage(ctx, &msg, msg.RoutingKey.Name)
		if err != nil {
			ignore[key] = true

			if err == utils.ErrorWaitPrevious {
				continue
			}

			msg.Attempts += 1
			maxRetryTimes := msg.RoutingKey.RetryTimes
			if maxRetryTimes <= 0 {
				maxRetryTimes = config.DefaultMaxRetryTimes
			}
			if msg.Attempts >= maxRetryTimes {
				msg.Status = models.InMessageStatusFailed
			}
		}

		err = i.inMsgRepo.UpdateMessage(&msg)
		if err != nil {
			logger.Errorf("Sent, failed to update status: %s, %s, %s, "+
				"error: %s", msg.RoutingKey.Name, msg.OriginModel, msg.OriginCode, err)
		}
	}
	logger.Info("[Retry Prev Message] Finish!")

	return nil
}

func (i *inService) CronArchivedMessages(ctx context.Context) error {
	query := schema.InMessageInStatusQueryParam{
		Status:      []string{models.InMessageStatusSuccess, models.InMessageStatusCanceled},
		CreatedDays: config.GetArchivedDays(),
		Page:        1,
		Limit:       config.GetArchivedSize(),
		Sort:        "_id",
	}
	messages, _, _ := i.inMsgRepo.GetMessagesInStatus(&query)
	if messages == nil {
		logger.Info("[Migrate In Message] Not found any message!")
		return nil
	}

	logger.Infof("[Migrate In Message] Found %d messages!", len(*messages))
	for _, msg := range *messages {
		err := i.inMsgRepo.ArchivedMessage(&msg)
		if err != nil {
			logger.Errorf("Failed to migrate in msg: %s, error: %s", msg.ID, err)
		}
	}
	logger.Info("[Migrate In Message] Finish!")

	return nil
}

func (i *inService) GetMessages(ctx context.Context, query *schema.InMessageQueryParam) (*[]models.InMessage, *paging.Paging, error) {
	result, pageInfo, err := i.inMsgRepo.GetMessages(query)
	if err != nil {
		return nil, nil, err
	}

	return result, pageInfo, errors.Success.New()
}

func (i *inService) UpdateMessage(ctx context.Context, id string,
	body *schema.InMessageBodyUpdateParam) (*models.InMessage, error) {
	message, err := i.inMsgRepo.GetMessageByID(id)
	if message == nil {
		return nil, utils.EMSGetDataNotfound.New()
	}

	var update models.InMessage
	copier.Copy(&update, &message)
	data, err := json.Marshal(body)
	if err != nil {
		return nil, utils.EMSSerializeError.New()
	}
	json.Unmarshal(data, &update)

	if reflect.DeepEqual(*message, update) {
		return message, errors.Success.New()
	}

	err = i.inMsgRepo.UpdateMessage(&update)
	if err != nil {
		return nil, err
	}

	return &update, errors.Success.New()
}

func (i *inService) handleMessage(ctx context.Context, message *models.InMessage, routingKey string) error {
	inRoutingKey, err := i.routingRepo.GetRoutingKey(
		&schema.RoutingKeyQueryParam{Name: routingKey},
	)
	if err != nil {
		err := errors.Wrap(err, fmt.Sprintf("cannot find routing key %s", routingKey))
		message.Status = models.InMessageStatusInvalid
		message.Logs = append(message.Logs, utils.ParseError(err))
		logger.Error(err)
		return err
	}
	message.RoutingKey = *inRoutingKey

	prevMsg, _ := i.getPrevMessage(message)
	if prevMsg != nil && prevMsg.Status != models.InMessageStatusSuccess && prevMsg.Status != models.InMessageStatusCanceled {
		message.Status = models.InMessageStatusWaitPrevMsg
		logger.Warn("Message must be wait for previous message!")
		return utils.ErrorWaitPrevious
	}

	res, err := i.externalService.CallAPI(ctx, message)
	if err != nil {
		message.Status = models.InMessageStatusWaitRetry
		if res != nil {
			message.Logs = append(message.Logs, res)
			return err
		}

		message.Logs = append(message.Logs, utils.ParseError(err))
		return err
	}

	message.Status = models.InMessageStatusSuccess
	message.Logs = append(message.Logs, res)

	return nil
}

func (i *inService) getPrevMessage(message *models.InMessage) (*models.InMessage, error) {
	query := schema.InMessagePreviousQueryParam{
		OriginModel: message.OriginModel,
		OriginCode:  message.OriginCode,
		PublishTime: message.PublishTime,
		Sort:        "-publish_time",
	}
	return i.inMsgRepo.GetPreviousMessage(&query)
}
