package services

import (
	"context"
	"encoding/json"
	"reflect"

	"github.com/jinzhu/copier"
	"github.com/spf13/viper"
	"transport/lib/errors"
	"transport/lib/utils/logger"
	"transport/lib/utils/paging"

	"transport/ems/app/models"
	"transport/ems/app/queue"
	"transport/ems/app/repositories"
	"transport/ems/app/schema"
	"transport/ems/config"
	"transport/ems/pkg/utils"
)

type OutService interface {
	Publish(ctx context.Context, body *schema.OutMessageBodyParam) error
	CronResend(ctx context.Context) error
	CronArchivedMessages(ctx context.Context) error
	GetMessages(ctx context.Context, query *schema.OutMessageQueryParam) (*[]models.OutMessage, *paging.Paging, error)
	UpdateMessage(ctx context.Context, id string, body *schema.OutMessageBodyUpdateParam) (*models.OutMessage, error)
	GetMessageByID(ctx context.Context, id string) (*models.OutMessage, error)
}

type outService struct {
	pub  queue.Publisher
	repo repositories.OutRepository
}

func NewOutService(pub queue.Publisher, repo repositories.OutRepository) OutService {
	return &outService{
		pub:  pub,
		repo: repo,
	}
}

func (o *outService) Publish(ctx context.Context, body *schema.OutMessageBodyParam) error {
	message, err := o.prepareOutMessage(body)
	if err != nil {
		logger.Errorf("Failed to prepare out message, %s", err)
		return err
	}

	err = o.pub.Publish(message, true)
	if err != nil {
		logger.Errorf("Failed to publish msg %s, %s", message.ID, err)
		return err
	}

	if viper.GetBool("publisher.ignore_store") {
		return nil
	}
	err = o.repo.AddMessage(message)
	if err != nil {
		logger.Errorf("Failed to create out msg %s", message.ID)
		return err
	}
	return nil
}

func (o *outService) CronResend(ctx context.Context) error {
	query := schema.OutMessageQueryParam{
		Status: models.OutMessageStatusWait,
		Page:   1,
		Limit:  config.CronResendOutMessageLimit,
		Sort:   "_id",
	}
	messages, _, _ := o.repo.GetMessages(&query)
	if messages == nil {
		logger.Info("[Resend Message] Not found any wait message!")
		return nil
	}

	logger.Infof("[Resend Message] Found %d wait messages!", len(*messages))
	for _, msg := range *messages {
		o.pub.Publish(&msg, true)
		o.repo.UpdateMessage(msg.ID, &msg)
	}
	logger.Info("[Resend Message] Finish!")

	return nil
}

func (o *outService) CronArchivedMessages(ctx context.Context) error {
	query := schema.OutMessageInStatusQueryParam{
		Status:      []string{models.OutMessageStatusSent, models.OutMessageStatusCanceled},
		CreatedDays: config.GetArchivedDays(),
		Page:        1,
		Limit:       config.CronArchivedMessageLimit,
		Sort:        "_id",
	}
	messages, _, _ := o.repo.GetMessagesInStatus(&query)
	if messages == nil {
		logger.Info("[Migrate Out Message] Not found any message!")
		return nil
	}

	logger.Infof("[Migrate Out Message] Found %d messages!", len(*messages))
	for _, msg := range *messages {
		err := o.repo.ArchivedMessage(&msg)
		if err != nil {
			logger.Errorf("Failed to migrate out msg: %s, error: %s", msg.ID, err)
		}
	}
	logger.Info("[Migrate Out Message] Finish!")

	return nil
}

func (o *outService) GetMessages(ctx context.Context, query *schema.OutMessageQueryParam) (*[]models.OutMessage, *paging.Paging, error) {
	result, pageInfo, err := o.repo.GetMessages(query)
	if err != nil {
		return nil, nil, err
	}

	return result, pageInfo, errors.Success.New()
}

func (o *outService) GetMessageByID(ctx context.Context, id string) (*models.OutMessage, error) {
	result, err := o.repo.GetMessageByID(id)
	if err != nil {
		return nil, err
	}

	return result, errors.Success.New()
}

func (o *outService) UpdateMessage(ctx context.Context, id string,
	body *schema.OutMessageBodyUpdateParam) (*models.OutMessage, error) {
	message, err := o.repo.GetMessageByID(id)
	if message == nil {
		return nil, utils.EMSGetDataNotfound.New()
	}

	var update models.OutMessage
	copier.Copy(&update, &message)
	data, err := json.Marshal(body)
	if err != nil {
		return nil, utils.EMSSerializeError.New()
	}
	json.Unmarshal(data, &update)

	if reflect.DeepEqual(*message, update) {
		return message, errors.Success.New()
	}

	err = o.repo.UpdateMessage(id, &update)
	if err != nil {
		return nil, err
	}

	return &update, errors.Success.New()
}

func (o *outService) prepareOutMessage(msgBody *schema.OutMessageBodyParam) (*models.OutMessage, error) {
	message := models.OutMessage{
		Status: models.OutMessageStatusWait,
	}

	err := copier.Copy(&message, msgBody)
	if err != nil {
		return nil, utils.EMSSerializeError.Wrap(err, "error")
	}

	return &message, nil
}
