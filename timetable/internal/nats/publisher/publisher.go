package natsPublisher

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/nats-io/nats.go"
	"timetable/internal/models"
)

var jsc nats.JetStreamContext

// Initialize NATS connection and JetStream context
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
		Name:     "ALERTS",
		Subjects: []string{"ALERTS.>"},
	})
	if err != nil {
		log.Fatal("Failed to create stream:", err)
	}
	log.Println("Publisher : NATS connected & stream initialized")
}

// Publish an alert to NATS
func PublishAlert(subject string, alert models.Alert) error {
	messageBytes, err := json.Marshal(alert)
	if err != nil {
		return err
	}

	// Publish the alert to the "alerts" subject
	pubAckFuture, err := jsc.PublishAsync("ALERTS." + subject, messageBytes)
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