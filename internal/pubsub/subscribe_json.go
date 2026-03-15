package pubsub

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type AckType int

const (
	Ack AckType = iota
	NackRequeue
	NackDiscard
)

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // an enum to represent "durable" or "transient"
	handler func(T) AckType,
) error {
	channel, queue, err := DeclareAndBind(conn, exchange, queueName, key, queueType)
	if err != nil {
		return err
	}

	deliveryChan, err := channel.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for message := range deliveryChan {
			var target T
			err := json.Unmarshal(message.Body, &target)
			if err != nil {
				log.Println(err)
				continue
			}
			ackType := handler(target)
			switch ackType {
			case Ack:
				log.Println("Ack called")
				message.Ack(false)
			case NackRequeue:
				log.Println("NackRequeue called")
				message.Nack(false, true)
			case NackDiscard:
				log.Println("NackDiscard called")
				message.Nack(false, false)
			}
		}
	}()
	return nil
}
