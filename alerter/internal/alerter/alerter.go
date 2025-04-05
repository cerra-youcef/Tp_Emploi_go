package alerter

import (
	"alerter/internal/models"
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// Connect to NATS and create a JetStream context
func ConnectToNATS() (jetstream.JetStream, *nats.Conn, error) {
	nc, err := nats.Connect(nats.DefaultURL) // Default: "nats://localhost:4222"
	if err != nil {
		return nil, nil, err
	}

	js, err := jetstream.New(nc)
	if err != nil {
		return nil, nil, err
	}
	return js, nc, nil
}

// Create or retrieve a Durable Consumer
func AlertConsumer(js jetstream.JetStream) (*jetstream.Consumer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get existing stream
	stream, err := js.Stream(ctx, "ALERTS")
	if err != nil {
		return nil, err
	}

	// Get or create a durable consumer
	consumer, err := stream.Consumer(ctx, "alertsConsumer")
	if err != nil {
		// Create if it doesn’t exist
		consumer, err = stream.CreateConsumer(ctx, jetstream.ConsumerConfig{
			Durable:     "alertsConsumer",
			Name:        "alertsConsumer",
			Description: "Consumes alerts from timetable",
		})
		if err != nil {
			return nil, err
		}
		log.Println("Created new consumer")
	} else {
		log.Println("Using existing consumer")
	}

	return &consumer, nil
}

func Consume(consumer jetstream.Consumer) error {

	cc, err := consumer.Consume(func(msg jetstream.Msg) {

		var alert models.AlertConsumer
		if err := json.Unmarshal(msg.Data(), &alert); err != nil {
			log.Println("Error decoding alert:", err)
			_ = msg.Nak()
			return
		}
		log.Println("-> alert type :", alert.Type)
		log.Println("event :", alert.Event)

		_ = msg.Ack()
	})

	<-cc.Closed()
	cc.Stop()

	return err
}
