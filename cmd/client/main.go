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
	fmt.Println("Starting Peril client...")

	connectionString := "amqp://guest:guest@127.0.0.1:5672/"
	connection, err := amqp.Dial(connectionString)
	if err != nil {
		log.Println(err)
		return
	}
	defer connection.Close()
	fmt.Println("Connected to RabbitMQ!")

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Println(err)
		return
	}
	username = fmt.Sprintf("%s.%s", routing.PauseKey, username)

	channel, queue, err := pubsub.DeclareAndBind(connection, routing.ExchangePerilDirect, username, routing.PauseKey, 1)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("connected to queue: %v on channel: %v", queue.Name, channel)

	gamestate := gamelogic.NewGameState(username)

	if 1 == 1 {
		input := gamelogic.GetInput()

		unitTypes := []string{"infantry", "cavalry", "artillery"}
		locations := []string{"americas", "europe", "africa", "asia", "antarctica", "australia"}

		found := false

		command := input[0]
		if command == "spawn" {
			if len(input) != 3 {
				fmt.Println("usage: spawn <type> <location>")
				return
			}
			location := input[1]
			unitType := input[2]

			for _, listLocation := range locations {
				if listLocation == unitType {
					found = true
				}
			}

			if !found {
				fmt.Println("usage: spawn <type> <location>")
				return
			}

			found = false

			for _, listUnitType := range unitTypes {
				if listUnitType == location {
					found = true
				}
			}

			if !found {
				fmt.Println("usage: spawn <type> <location>")
				return
			}

			err = gamestate.CommandSpawn(input)
			if err != nil {
				log.Println(err)
				return
			}
		} else if command == "move" {
			if len(input) != 3 {
				fmt.Println("usage: move <from> <to>")
				return
			}
		}
	}
	// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
}
