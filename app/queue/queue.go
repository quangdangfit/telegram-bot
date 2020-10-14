package queue

import (
	"time"

	"github.com/streadway/amqp"
	"transport/lib/utils/logger"

	"telegram-bot/pkg/utils"
)

const (
	RecoverIntervalTime = 6 * 60
	TimeoutRetry        = 3
	WaitTimeReconnect   = 5
)

type MessageQueue struct {
	config          *AMQPConfig
	connection      *amqp.Connection
	channel         *amqp.Channel
	errorChan       chan *amqp.Error
	failedChan      chan amqp.Delivery
	isClosed        bool
	channelIsClosed bool
}

func (mq *MessageQueue) newConnection() error {
	conn, err := amqp.Dial(mq.config.AMQPUrl)
	for err != nil {
		logger.Error("Failed to create new connection to AMQP: ", err)

		logger.Infof("Sleep %d seconds to reconnect", WaitTimeReconnect)
		time.Sleep(WaitTimeReconnect * time.Second)
		conn, err = amqp.Dial(mq.config.AMQPUrl)
	}
	mq.connection = conn

	return nil
}

func (mq *MessageQueue) closeConnection() error {
	if mq.isClosed {
		return nil
	}
	mq.closeChannel()

	if mq.connection != nil {
		if err := mq.connection.Close(); err != nil {
			return utils.EMSAMQPCloseConnectionError.Wrap(err, "error")
		}
		mq.connection = nil
	}

	mq.isClosed = true
	return nil
}

func (mq *MessageQueue) newChannel() (*amqp.Channel, error) {
	mq.ensureConnection()

	if mq.connection == nil || mq.connection.IsClosed() {
		logger.Error("Connection is not open, cannot create new channel")
		return nil, utils.EMSAMQPConnectionClosedError.Newm("connection is closed")
	}

	channel, err := mq.connection.Channel()
	if err != nil {
		logger.Error("Failed to new channel: ", err)
		return nil, utils.EMSAMQPOpenChannelError.Wrap(err, "error")
	}
	mq.channel = channel
	mq.channelIsClosed = false
	logger.Info("New channel successfully")
	return channel, nil
}

func (mq *MessageQueue) ensureConnection() (err error) {
	if mq.connection == nil || mq.connection.IsClosed() {
		err = mq.newConnection()
		if err != nil {
			return err
		}
	}
	return nil
}

func (mq *MessageQueue) closeChannel() error {
	if mq.isClosed || mq.channelIsClosed {
		return nil
	}
	logger.Info("Close channel")
	if mq.channel != nil {
		_ = mq.channel.Close()
		mq.channel = nil
		mq.channelIsClosed = true
	}

	return nil
}

func (mq *MessageQueue) declareExchange() error {
	mq.newChannel()
	defer mq.closeChannel()

	if mq.ChanelIsClosed() {
		logger.Error("Channel is not open, cannot declare exchange")
		return utils.EMSAMQPChannelClosedError.New()
	}

	if err := mq.channel.ExchangeDeclare(
		mq.config.ExchangeName, // name
		mq.config.ExchangeType, // type
		true,                   // durable
		false,                  // auto-deleted
		false,                  // internal
		false,                  // noWait
		nil,                    // arguments
	); err != nil {
		logger.Error("Failed to declare exchange: ", err)
		return utils.EMSAMQPExchangeDeclareError.Wrap(err, "error")
	}

	logger.Info("Declared exchange: ", mq.config.ExchangeName)
	return nil
}

func (mq *MessageQueue) declareQueue() error {
	mq.newChannel()
	defer mq.closeChannel()

	if mq.ChanelIsClosed() {
		logger.Error("Channel is not open, cannot declare exchange")
		return utils.EMSAMQPChannelClosedError.New()
	}

	if _, err := mq.channel.QueueDeclare(
		mq.config.QueueName,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		logger.Error("Failed to declare queue: ", err)
		return utils.EMSAMQPQueueDeclareError.Wrap(err, "error")
	}

	logger.Info("Declared queue: ", mq.config.QueueName)
	return nil
}

func (mq *MessageQueue) bindQueue(routingKey string) error {
	if err := mq.channel.QueueBind(
		mq.config.QueueName,    // name
		routingKey,             // key
		mq.config.ExchangeName, // exchange
		false,                  //noWait
		nil,                    // args
	); err != nil {
		logger.Error("Failed to bind queue: ", err)
		return utils.EMSAMQPQueueBindError.Wrap(err, "error")
	}
	return nil
}

func (mq *MessageQueue) ChanelIsClosed() bool {
	if mq.channel == nil || mq.channelIsClosed {
		return true
	}
	return false
}
