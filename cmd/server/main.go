package main

import (
	"encoding/json"
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

	pausedPlayingState, err := json.Marshal(routing.PlayingState{IsPaused: true})
	if err != nil {
		log.Println(err)
	}

	resumePlayingState, err := json.Marshal(routing.PlayingState{IsPaused: false})
	if err != nil {
		log.Println(err)
	}

	if 1 == 1 {
		input := gamelogic.GetInput()

		command := input[0]

		if command == "pause" {
			log.Println("Pausing game...")
			err = pubsub.PublishJSON(mainChannel, routing.ExchangePerilDirect, routing.PauseKey, pausedPlayingState)
		} else if command == "resume" {
			log.Println("Resuming game...")
			err = pubsub.PublishJSON(mainChannel, routing.ExchangePerilDirect, routing.PauseKey, resumePlayingState)
		} else if command == "quit" {
			log.Println("Quitting...")
			return
		} else {
			log.Println("Invalid command.")
		}
	}

	// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
}
