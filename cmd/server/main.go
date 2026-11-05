package main

import (
	"log"
	"os"
	"os/signal"

	amqp "github.com/rabbitmq/amqp091-go"
)

const connURL = "amqp://guest:guest@localhost:5672/"

func main() {
	conn, err := amqp.Dial(connURL)
	if err != nil {
		log.Fatalf("Failed to connect to rabbitmq: %v", err)
	}

	defer func() {
		if err := conn.Close(); err != nil {
			log.Fatalf("Failed to close rabbitmq connection: %v", err)
		}
	}()

	log.Println("Starting Peril server...")

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)
	<-signalCh

	log.Println("Server closed")
}
