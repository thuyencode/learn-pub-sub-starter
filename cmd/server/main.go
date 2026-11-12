package main

import (
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

const connURL = "amqp://guest:guest@localhost:5672/"

func main() {
	log.Println("Starting Peril server...")

	conn, err := amqp.Dial(connURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	defer func() {
		if err := conn.Close(); err != nil {
			log.Fatalf("Failed to close RabbitMQ connection: %v", err)
		}
	}()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to create a AMQP channel: %v", err)
	}

	_, queue, err := pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilTopic,
		routing.GameLogSlug,
		routing.GameLogSlug+".*",
		pubsub.DurableQueue,
	)
	if err != nil {
		log.Fatalf("Failed to declare and bind a queue: %v", err)
	}
	log.Printf("Queue %s declared and bound", queue.Name)

	gamelogic.PrintServerHelp()

loop:
	for {
		input := gamelogic.GetInput()
		if len(input) == 0 {
			continue
		}

		switch input[0] {
		case "pause":
			log.Println("Sending a pause message...")
			if err := pubsub.PublishJSON(
				ch,
				routing.ExchangePerilDirect,
				routing.PauseKey,
				routing.PlayingState{IsPaused: true},
			); err != nil {
				log.Fatalf("Failed to publish a message: %v", err)
			}
		case "resume":
			log.Println("Sending a resume message...")
			if err := pubsub.PublishJSON(
				ch,
				routing.ExchangePerilDirect,
				routing.PauseKey,
				routing.PlayingState{IsPaused: false},
			); err != nil {
				log.Fatalf("Failed to publish a message: %v", err)
			}
		case "quit":
			break loop
		case "help":
			gamelogic.PrintServerHelp()
		default:
			log.Println("No such command")
		}
	}

	log.Println("Server closed")
}
