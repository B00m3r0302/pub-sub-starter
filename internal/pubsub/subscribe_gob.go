package pubsub

import (
	"bytes"
	"encoding/gob"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SubscribeGob[T any](
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

	err = channel.Qos(10, 0, false)

	deliveryChan, err := channel.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for message := range deliveryChan {
			var target T
			buf := bytes.NewReader(message.Body)
			dec := gob.NewDecoder(buf)
			err := dec.Decode(&target)
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
