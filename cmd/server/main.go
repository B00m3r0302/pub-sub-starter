package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")

	connectionString := "amqp://guest:guest@127.0.0.1:5672/"
	connection, err := amqp.Dial(connectionString)
	if err != nil {
		log.Println(err)
		return
	}
	defer connection.Close()
	fmt.Println("Connected to RabbitMQ!")

	// See commands user can use
	gamelogic.PrintServerHelp()

	// Make a new channel
	mainChannel, err := connection.Channel()
	if err != nil {
		log.Println(err)
		return
	}
	defer mainChannel.Close()

	topicChannel, topicQueue, err := pubsub.DeclareAndBind(connection, routing.ExchangePerilTopic, routing.GameLogSlug, routing.LogsKey, 0)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("connected to queue: %v on channel: %v", topicQueue.Name, topicChannel)

	for {
		input := gamelogic.GetInput()

		command := input[0]
		switch command {
		case "pause":
			log.Println("Pausing game...")
			err = pubsub.PublishJSON(mainChannel, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: true})
		case "resume":
			log.Println("Resuming game...")
			err = pubsub.PublishJSON(mainChannel, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: false})
		case "quit":
			log.Println("Quitting...")
			return
		default:
			log.Println("Invalid command.")
		}
	}

	// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
}
