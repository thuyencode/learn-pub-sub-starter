package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

const connURL = "amqp://guest:guest@localhost:5672/"

func main() {
	fmt.Println("Starting Peril client...")

	conn, err := amqp.Dial(connURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	defer func() {
		if err := conn.Close(); err != nil {
			log.Fatalf("Failed to close RabbitMQ connection: %v", err)
		}
	}()

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("Failed to get username: %v", err)
	}

	if _, _, err := pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilDirect,
		"pause."+username,
		routing.PauseKey,
		pubsub.TransientQueue,
	); err != nil {
		log.Fatalf("Failed to declare and bind a queue: %v", err)
	}

	gameState := gamelogic.NewGameState(username)

loop:
	for {
		input := gamelogic.GetInput()
		if len(input) == 0 {
			continue
		}

		switch input[0] {
		case "spawn":
			if len(input[0:]) < 3 {
				fmt.Println(`Not enough arguments. See "help"`)
			} else if err := gameState.CommandSpawn(input[0:]); err != nil {
				fmt.Println(err)
			}
		case "move":
			if len(input[0:]) < 3 {
				fmt.Println(`Not enough arguments. See "help"`)
			} else {
				move, err := gameState.CommandMove(input[0:])
				if err != nil {
					fmt.Println(err)
				}
				fmt.Printf("Moved to %s\n", move.ToLocation)
			}
		case "status":
			gameState.CommandStatus()
		case "spam":
			fmt.Println("Spamming not allowed yet!")
		case "help":
			gamelogic.PrintClientHelp()
		case "quit":
			gamelogic.PrintQuit()
			break loop
		default:
			fmt.Println("No such command")
		}
	}
}
