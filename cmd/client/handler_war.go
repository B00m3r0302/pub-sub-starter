package main

import (
	"fmt"
	"log"
	"time"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func publishGameLog(ch *amqp.Channel, username, message string) error {
	key := fmt.Sprintf("%s.%s", routing.GameLogSlug, username)
	Exchange := routing.ExchangePerilTopic
	GameLog := routing.GameLog{
		CurrentTime: time.Now(),
		Message:     message,
		Username:    username,
	}

	err := pubsub.PublishGob(ch, Exchange, key, GameLog)
	if err != nil {
		return err
	}
	return nil
}

func handlerWar(gs *gamelogic.GameState, ch *amqp.Channel) func(gamelogic.RecognitionOfWar) pubsub.AckType {
	return func(dw gamelogic.RecognitionOfWar) pubsub.AckType {
		defer fmt.Print("> ")
		outcome, winner, loser := gs.HandleWar(dw)

		switch outcome {
		case gamelogic.WarOutcomeNotInvolved:
			return pubsub.NackRequeue
		case gamelogic.WarOutcomeNoUnits:
			return pubsub.NackDiscard
		case gamelogic.WarOutcomeYouWon:
			msg := fmt.Sprintf("%s won a war against %s", winner, loser)
			log.Println(msg)
			err := publishGameLog(ch, gs.GetUsername(), msg)
			if err != nil {
				return pubsub.NackRequeue
			}
			return pubsub.Ack
		case gamelogic.WarOutcomeOpponentWon:
			msg := fmt.Sprintf("%s won a war against %s", winner, loser)
			log.Println(msg)
			err := publishGameLog(ch, gs.GetUsername(), msg)
			if err != nil {
				return pubsub.NackRequeue
			}
			return pubsub.Ack
		case gamelogic.WarOutcomeDraw:
			msg := fmt.Sprintf("A war between %s and %s resulted in a draw", winner, loser)
			log.Println(msg)
			err := publishGameLog(ch, gs.GetUsername(), msg)
			if err != nil {
				return pubsub.NackRequeue
			}
			return pubsub.Ack
		}
		log.Println("Some weird outcome happened")
		return pubsub.NackDiscard
	}
}
