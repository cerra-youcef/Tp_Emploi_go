package natsConsumer

import (
	"context"
	"encoding/json"
	"log"
	"time"
	"timetable/internal/models"
	"timetable/internal/services/Events"

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
func EventConsumer(js jetstream.JetStream) (*jetstream.Consumer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get existing stream
	stream, err := js.Stream(ctx, "EVENTS")
	if err != nil {
		return nil, err
	}

	// Get or create a durable consumer
	consumer, err := stream.Consumer(ctx, "timetable_consumer")
	if err != nil {
		// Create if it doesn’t exist
		consumer, err = stream.CreateConsumer(ctx, jetstream.ConsumerConfig{
			Durable:     "TimetableConsumer",
			Name:        "TimetableConsumer",
			Description: "Consumes events for timetable",
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

// Consume messages from NATS and store in DB
func Consume(ctx context.Context, consumer jetstream.Consumer) error {
	cc, err := consumer.Consume(func(msg jetstream.Msg) {
		
		var event models.Event
		if err := json.Unmarshal(msg.Data(), &event); err != nil {
			log.Println("Error decoding event:", err)
			_ = msg.Nak()
			return
		}
		log.Printf("Received event: ID: %s, Resource: %d, Name: %s\n", event.ID,event.Resource,event.Name )

		// Store in database if it doesn't exist
		existingEvent, err := Events.GetEventByID(event.ID)
		if err != nil {
			log.Println("Error checking event existence:", err)
			_ = msg.Nak()
			return
		}

		if existingEvent == nil {
			err := Events.CreateEvent(&event)
			if err != nil {
				log.Println("Error storing event:", err)
			} else {
				log.Println("New event stored:", event.ID)
			}
		} else {
			log.Println("Event already exists, skipping:", event.ID)
		}
		_ = msg.Ack()
	})

	// Keep the consumer running
	<-cc.Closed()
	cc.Stop()

	return err
}