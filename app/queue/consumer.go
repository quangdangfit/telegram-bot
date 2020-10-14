package queue

import (
	"encoding/json"
	"sync/atomic"
	"time"

	"github.com/jinzhu/copier"
	"github.com/manucorporat/try"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
	"transport/lib/utils/logger"

	"transport/ems/app/models"
	"transport/ems/pkg/utils"
)

type Consumer interface {
	Consume() (chan *models.InMessage, error)
}

type consumer struct {
	MessageQueue
	done            chan error
	consumerTag     string // Name that consumer identifies itself to the server
	lastRecoverTime int64
	//track service current status
	currentStatus atomic.Value

	msgChan chan *models.InMessage
}

func NewConsumer() Consumer {
	var cons = consumer{
		done:            make(chan error),
		lastRecoverTime: time.Now().Unix(),
	}

	cons.config = &AMQPConfig{
		AMQPUrl:      viper.GetString("ts_rabbit.amqp"),
		QueueName:    viper.GetString("ts_rabbit.queue_name"),
		ExchangeName: viper.GetString("ts_rabbit.exchange_name"),
		ExchangeType: viper.GetString("ts_rabbit.exchange_type"),
	}
	err := cons.newConnection()
	if err != nil {
		logger.Error("Consumer failed to open new connection!")
	}

	err = cons.declareQueue()
	if err != nil {
		logger.Error("Consumer failed to declare queue!")
	}

	cons.msgChan = make(chan *models.InMessage)

	return &cons
}

func (cons *consumer) Consume() (chan *models.InMessage, error) {
	cons.ensureConnection()
	cons.newChannel()
	deliveries, err := cons.subscribe()
	if err != nil {
		return nil, utils.EMSAMQPConsumerSubscribeError.Wrap(err, "error")
	}
	go cons.startConsuming(deliveries)

	// retry consume failed messages
	go cons.startConsuming(cons.failedChan)

	return cons.msgChan, nil
}

func (cons *consumer) parseMessageFromDelivery(msg amqp.Delivery) (*models.InMessage, error) {
	var payload interface{}
	json.Unmarshal(msg.Body, &payload)

	var headers models.Headers
	data, err := json.Marshal(msg.Headers)
	if err != nil {
		return nil, utils.EMSSerializeError.Wrap(err, "error")
	}
	json.Unmarshal(data, &headers)

	message := models.InMessage{
		Payload: payload,
	}
	copier.Copy(&message, &headers)
	message.RoutingKey.Name = msg.RoutingKey

	return &message, nil
}

func (cons *consumer) reconnect(retryTime int) (<-chan amqp.Delivery, error) {
	cons.closeConnection()
	time.Sleep(time.Duration(TimeoutRetry) * time.Second)
	logger.Info("Try reConnect with times:", retryTime)

	cons.ensureConnection()

	deliveries, err := cons.subscribe()
	if err != nil {
		return deliveries, utils.EMSAMQPConsumerSubscribeError.Wrap(err, "error")
	}
	return deliveries, nil
}

// subscribe sets the queue that will be listened to for this connection
func (cons *consumer) subscribe() (<-chan amqp.Delivery, error) {
	err := cons.channel.Qos(50, 0, false)
	if err != nil {
		logger.Error("Error setting qos: ", err)
		return nil, utils.EMSAMQPConsumerQosError.Wrap(err, "error")
	}

	logger.Info("Queue bound to Exchange, starting Consume consumer tag:", cons.consumerTag)

	deliveries, err := cons.channel.Consume(
		cons.config.QueueName, // name
		cons.consumerTag,      // consumerTag,
		false,                 // noAck
		false,                 // exclusive
		false,                 // noLocal
		false,                 // noWait
		nil,                   // arguments
	)
	if err != nil {
		logger.Error("Failed to consume queue: ", err)
		return nil, utils.EMSAMQPConsumerConsumeError.Wrap(err, "error")
	}
	return deliveries, nil
}

func (cons *consumer) startConsuming(deliveries <-chan amqp.Delivery) {
	for msg := range deliveries {
		logger.Info("Enter deliver message: ", msg.RoutingKey)
		ret := false
		try.This(func() {
			cons.pushToMsgChan(msg)
		}).Finally(func() {
			if ret == true {
				msg.Ack(false)
				currentTime := time.Now().Unix()
				if currentTime-cons.lastRecoverTime > RecoverIntervalTime &&
					!cons.currentStatus.Load().(bool) {

					logger.Info("Try to Recover Unack Messages!")
					cons.currentStatus.Store(true)
					cons.lastRecoverTime = currentTime
					cons.channel.Recover(true)
				}

			} else {
				// this really a litter dangerous. if the worker is panic very
				//quickly, it will ddos our sentry server......plz,
				//add [retry-ttl] in header.
				//msg.Nack(false, true)
				msg.Reject(false)
				//c.currentStatus.Store(true)
			}
		}).Catch(func(e try.E) {
			logger.Error(e)
		})
	}
}

func (cons *consumer) pushToMsgChan(msg amqp.Delivery) {
	message, err := cons.parseMessageFromDelivery(msg)
	if err != nil {
		cons.failedChan <- msg
		logger.Error("Failed to parse message: ", err)
	}

	cons.msgChan <- message
}
