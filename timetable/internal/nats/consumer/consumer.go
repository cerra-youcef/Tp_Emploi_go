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
	consumer, err := stream.Consumer(ctx, "TimetableConsumer")
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
func Consume(consumer jetstream.Consumer) error {

	var receivedEvents []models.Event // Store received events
	cc, err := consumer.Consume(func(msg jetstream.Msg) {

		subject := msg.Subject()

		//if all events are received we can check for removed ones
		if subject == "EVENTS.end" {
			if err := Events.DeleteRemovedEvents(receivedEvents); err != nil {
				log.Println("Error Deleting Event and publishing alert : ", err)
			}
			_ = msg.Ack()
			return
		}

		var event models.Event
		if err := json.Unmarshal(msg.Data(), &event); err != nil {
			log.Println("Error decoding event:", err)
			_ = msg.Nak()
			return
		}
		receivedEvents = append(receivedEvents, event) // Track received events

		existingEvent, err := Events.GetEventByUID(event.UID)
		if err != nil {
			log.Println("Error checking event existence:", err)
			_ = msg.Nak()
			return
		}
		//create event if it doesnt exist or update it otherwise & publish alert
		if existingEvent == nil {
			if err := Events.CreateAndNotifyEvent(event); err != nil {
				log.Println("Error Creating Event and publishing alert : ", err)
			}
		} else {
			if err := Events.UpdateAndNotifyEvent(*existingEvent, event); err != nil {
				log.Println("Error Updating Event and publishing alert : ", err)
			}
		}
		_ = msg.Ack()
	})

	<-cc.Closed()
	cc.Stop()

	return err
}
