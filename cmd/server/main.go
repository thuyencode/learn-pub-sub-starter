package main

import (
	"log"
	"os"
	"os/signal"

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

	if err := pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: true}); err != nil {
		log.Fatalf("Failed to publish a message: %v", err)
	}

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)
	<-signalCh

	log.Println("Server closed")
}
