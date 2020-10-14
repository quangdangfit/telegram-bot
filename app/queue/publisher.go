package queue

import (
	"encoding/json"
	"time"

	"github.com/spf13/viper"
	"github.com/streadway/amqp"
	"transport/lib/utils/logger"

	"telegram-bot/app/models"
	"telegram-bot/pkg/utils"
)

type Publisher interface {
	Publish(message *models.OutMessage, reliable bool) error
}

type publisher struct {
	MessageQueue
}

func NewPublisher() Publisher {
	var pub publisher
	pub.config = &AMQPConfig{
		AMQPUrl:      viper.GetString("ts_rabbit.amqp"),
		ExchangeName: viper.GetString("ts_rabbit.exchange_name"),
		ExchangeType: viper.GetString("ts_rabbit.exchange_type"),
		QueueName:    viper.GetString("ts_rabbit.queue_name"),
	}
	err := pub.newConnection()
	if err != nil {
		logger.Error("Publisher failed to create new connection ", err)
	}

	err = pub.declareExchange()
	if err != nil {
		logger.Error("Publisher failed to declare exchange ", err)
	}

	if viper.GetString("ts_service.mode") == "publisher" {
		err = pub.declareQueue()
		if err != nil {
			logger.Error("Publisher failed to declare queue ", err)
		}
	}

	return &pub
}

func (pub *publisher) Publish(message *models.OutMessage, reliable bool) error {
	// New channel and close after publish
	pub.ensureConnection()
	channel, _ := pub.connection.Channel()
	defer channel.Close()

	// Reliable publisher confirms require confirm.select support from the connection.
	var confirms chan amqp.Confirmation
	if reliable {
		if err := channel.Confirm(false); err != nil {
			logger.Error("Channel could not be put into confirm mode ", err)
			return utils.EMSAMQPublisherConfirmError.Wrap(err, "error")
		}
		confirms = channel.NotifyPublish(make(chan amqp.Confirmation, 1))
	}

	payload, err := json.Marshal(message.Payload)
	if err != nil {
		return utils.EMSSerializeError.Wrap(err, "error")
	}

	headers, err := pub.prepareHeaders(message)
	if err != nil {
		return err
	}

	err = channel.Publish(
		pub.config.ExchangeName, // publish to an exchange
		message.RoutingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			Headers:         headers,
			ContentType:     "application/json",
			ContentEncoding: "",
			Body:            payload,
			DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
			Priority:        0,              // 0-9
			Timestamp:       time.Now(),
			// a bunch of application/implementation-specific fields
		},
	)

	if err != nil {
		message.Status = models.OutMessageStatusFailed
		message.Logs = append(message.Logs, utils.ParseError(err))
		logger.Error("Failed to publish message ", err)
		return utils.EMSAMQPublisherPublishError.Wrap(err, "error")
	}

	if confirms != nil {
		pub.confirmOne(message, confirms)
	}

	return nil
}

func (pub *publisher) prepareHeaders(message *models.OutMessage) (amqp.Table, error) {
	strQuery, err := utils.URLEncode(message.Query)
	if err != nil {
		return nil, utils.EMSParseQueryError.New()
	}

	var headers = amqp.Table{
		"origin_code":  message.OriginCode,
		"origin_model": message.OriginModel,
		"query":        strQuery,
		"external_id":  message.ExternalID,
		"publish_time": time.Now().UTC().Format(time.RFC3339Nano),
	}

	if message.Headers != nil {
		headersByte, err := json.Marshal(message.Headers)
		if err != nil {
			return nil, utils.EMSSerializeError.New()
		}
		json.Unmarshal(headersByte, &headers)
	}

	return headers, nil
}

func (pub *publisher) confirmOne(message *models.OutMessage,
	confirms <-chan amqp.Confirmation) bool {

	confirmed := <-confirms
	if confirmed.Ack {
		logger.Info("Confirmed delivery with delivery tag: ", confirmed.DeliveryTag)

		message.Status = models.OutMessageStatusSent
		return true
	}

	logger.Info("Failed delivery of delivery tag: ", confirmed.DeliveryTag)

	message.Status = models.OutMessageStatusSentWait
	return false
}
