package pubsub

import amqp "github.com/rabbitmq/amqp091-go"

type SimpleQueueType int

const (
	Durable SimpleQueueType = iota
	Transient
)

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // SimpleQueueType is an "enum" type I made to represent "durable" or "transient"
) (*amqp.Channel, amqp.Queue, error) {
	channel, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, err
	}
	defer channel.Close()

	var durable bool
	if queueType == 0 {
		durable = true
	}

	if queueType == 1 {
		durable = false
	}

	var autoDelete bool
	if queueType == 0 {
		autoDelete = false
	}
	if queueType == 1 {
		autoDelete = true
	}

	var exclusive bool
	if queueType == 0 {
		exclusive = false
	}
	if queueType == 1 {
		exclusive = true
	}

	queue, err := channel.QueueDeclare(queueName, durable, autoDelete, exclusive, false, nil)

	err = channel.QueueBind(queue.Name, key, exchange, false, nil)
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	return channel, queue, nil
}
