package pubsub

import (
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // an enum to represent "durable" or "transient"
	handler func(T),
) error {
	ch, queue, err := DeclareAndBind(conn, exchange, queueName, key, queueType)
	if err != nil {
		return fmt.Errorf("failed to declare and bind a queue: %w", err)
	}

	deliveryCh, err := ch.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to consume queue %s: %w", queue.Name, err)
	}

	go func() {
		for delivery := range deliveryCh {
			var message T
			if err := json.Unmarshal(delivery.Body, &message); err != nil {
				log.Printf("failed to unmarshal message from queue %s: %v\n", queue.Name, err)
				continue
			}

			handler(message)

			if err := delivery.Ack(false); err != nil {
				log.Printf("failed to ack a delivery from queue %s: %v\n", queue.Name, err)
			}
		}
	}()

	return nil
}
