package nats

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/nats-io/nats.go"
)

var jsc nats.JetStreamContext

func InitNATS() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal("Failed to connect to NATS:", err)
	}
	
	// Get JetStream context
	jsc, err = nc.JetStream()
	if err != nil {
		log.Fatal("Failed to get JetStream context:", err)
	}

	// Initialize stream
	_, err = jsc.AddStream(&nats.StreamConfig{
		Name:     "EVENTS",
		Subjects: []string{"EVENTS.>"},
	})
	if err != nil {
		log.Fatal("Failed to create stream:", err)
	}

	log.Println("NATS connected & stream initialized")
}

// Publish an event to NATS
func PublishEvent(subject string, eventData interface{}) error {
	messageBytes, err := json.Marshal(eventData)
	if err != nil {
		return err
	}

	pubAckFuture, err := jsc.PublishAsync(subject, messageBytes)
	if err != nil {
		return err
	}

	select {
	case <-pubAckFuture.Ok():
		return nil
	case <-pubAckFuture.Err():
		return errors.New(string(pubAckFuture.Msg().Data))
	}
}
