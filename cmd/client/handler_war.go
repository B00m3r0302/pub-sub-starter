package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func publishGameLog(ch *amqp.Channel, exchange, username string, gl routing.GameLog) error {
	key := fmt.Sprintf("%s.%s", routing.games)
}

func handlerWar(gs *gamelogic.GameState) func(gamelogic.RecognitionOfWar) pubsub.AckType {
	return func(dw gamelogic.RecognitionOfWar) pubsub.AckType {
		defer fmt.Print("> ")
		outcome, winner, loser := gs.HandleWar(dw)

		switch outcome {
		case gamelogic.WarOutcomeNotInvolved:
			return pubsub.NackRequeue
		case gamelogic.WarOutcomeNoUnits:
			return pubsub.NackDiscard
		case gamelogic.WarOutcomeYouWon:
			log.Printf("%s won a war against %s", winner, loser)
			return pubsub.Ack
		case gamelogic.WarOutcomeOpponentWon:
			log.Printf("%s won a war against %s", winner, loser)
			return pubsub.Ack
		case gamelogic.WarOutcomeDraw:
			log.Printf("A war between %s and %s resulted in a draw", winner, loser)
			return pubsub.Ack
		}
		log.Println("Some weird outcome happened")
		return pubsub.NackDiscard
	}
}
