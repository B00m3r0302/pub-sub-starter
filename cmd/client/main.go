package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

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
	pauseQueueName := fmt.Sprintf("%s.%s", routing.PauseKey, username)
	movesQueueName := fmt.Sprintf("%s.%s", routing.ArmyMovesPrefix, username)
	gamestate := gamelogic.NewGameState(username)

	pauseHandler := handlerPause(gamestate)
	err = pubsub.SubscribeJSON(connection, routing.ExchangePerilDirect, pauseQueueName, routing.PauseKey, 1, pauseHandler)
	if err != nil {
		log.Println(err)
		return
	}

	publishCh, err := connection.Channel()
	if err != nil {
		log.Println(err)
	}

	movesHandler := handlerMove(gamestate, publishCh)
	err = pubsub.SubscribeJSON(connection, routing.ExchangePerilTopic, movesQueueName, "army_moves.*", 1, movesHandler)
	if err != nil {
		log.Println(err)
		return
	}

	warHandler := handlerWar(gamestate, publishCh)
	err = pubsub.SubscribeJSON(connection, routing.ExchangePerilTopic, routing.WarRecognitionsPrefix, routing.WarRecognitionsPrefix+".*", 0, warHandler)
	if err != nil {
		log.Println(err)
		return
	}
	for {
		input := gamelogic.GetInput()

		unitTypes := []string{"infantry", "cavalry", "artillery"}
		locations := []string{"americas", "europe", "africa", "asia", "antarctica", "australia"}

		found := false

		if len(input) == 0 {
			continue
		}

		command := input[0]
		switch command {
		case "spawn":
			if len(input) != 3 {
				fmt.Println("usage: spawn <type> <location>")
				continue
			}
			location := input[1]
			unitType := input[2]

			for _, listLocation := range locations {
				if listLocation == location {
					found = true
				}
			}

			if !found {
				fmt.Println("usage: spawn <location> <type>")
				continue
			}

			found = false

			for _, listUnitType := range unitTypes {
				if listUnitType == unitType {
					found = true
				}
			}

			if !found {
				fmt.Println("usage: spawn <type> <location>")
				continue
			}

			err = gamestate.CommandSpawn(input)
			if err != nil {
				log.Println(err)
				continue
			}
		case "move":
			if len(input) != 3 {
				fmt.Println("usage: move <from> <to>")
				continue
			}

			move, err := gamestate.CommandMove(input)
			if err != nil {
				log.Println(err)
				continue
			}

			moveUsername := fmt.Sprintf("%s.%s", routing.ArmyMovesPrefix, move.Player.Username)
			pubsub.PublishJSON(publishCh, routing.ExchangePerilTopic, moveUsername, move)
			log.Println("move was published to RabbitMQ!", input)
		case "status":
			gamestate.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			if len(input) != 2 {
				fmt.Println("usage: spawn <spam count>")
				continue
			}
			stripped := strings.TrimSpace(input[1])
			number, err := strconv.Atoi(stripped)
			if err != nil {
				fmt.Println("conversion of spam number to int failed", err)
				continue
			}
			for i := 0; i < number; i++ {
				msg := gamelogic.GetMaliciousLog()
				err = publishGameLog(publishCh, username, msg)
				if err != nil {
					fmt.Printf("Error publishing malicious log: %s\n", err)
				}
			}
			fmt.Printf("Published %v malicious logs\n", number)
		case "quit":
			gamelogic.PrintQuit()
			return
		default:
			fmt.Println("Invalid command.")
		}
	}
}
