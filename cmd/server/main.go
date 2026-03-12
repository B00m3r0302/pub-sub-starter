package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")

	connectionString := "amqp://guest:guest@127.0.0.1:5672/"
	connection, err := amqp.Dial(connectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()
	fmt.Println("Connected to RabbitMQ!")

	// Make a new channel
	mainChannel, err := connection.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer mainChannel.Close()

	pausedPlayingState, err := json.Marshal(routing.PlayingState{IsPaused: true})
	if err != nil {
		log.Fatal(err)
	}

	err = pubsub.PublishJSON(mainChannel, routing.ExchangePerilDirect, routing.PauseKey, pausedPlayingState)
	// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
}
