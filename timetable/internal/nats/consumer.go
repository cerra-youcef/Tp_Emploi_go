package natsConsumer

import (
	"context"
	"encoding/json"
	"log"
	"time"
	"timetable/internal/models"
	"timetable/internal/services/Events"
	"timetable/internal/helpers"

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
func Consume(consumer jetstream.Consumer) error {
	//var receivedEvents []models.Event // Store received events

	cc, err := consumer.Consume(func(msg jetstream.Msg) {
		
		var event models.Event

		if err := json.Unmarshal(msg.Data(), &event); err != nil {
			log.Println("Error decoding event:", err)
			_ = msg.Nak()
			return
		}
		//receivedEvents = append(receivedEvents, event) // Track received events

		//log.Printf("Received event: uid : %s Name: %s\n",event.UID,event.Name )

		// Store in database if it doesn't exist
		existingEvent, err := Events.GetEventByUID(event.UID)
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
					log.Println("New event stored:", event.UID)
					helpers.PrintAlert("event.created", event)
			}
		} else {
			// Compare changes and update (event exist)
			changes := helpers.DetectChanges(*existingEvent, event)
			if len(changes) > 0 {
				log.Println("Event updated, publishing alert:")
				helpers.PrintAlert("event.updated", event, changes)
				err := Events.UpdateEvent(&event)
				if err != nil {
					log.Println("Error updating event:", err)
				}
			}
		}
		//helpers.DetectDeletedEvents(receivedEvents) // Check for deleted events once all events are received
		_ = msg.Ack()
	})


	<-cc.Closed() // Wait for the consumer to stop receiving messages
	cc.Stop()

	return err
}


// Detect and handle deleted events
func DetectDeletedEvents(receivedEvents []models.Event) {

	if len(receivedEvents) == 0 {
		log.Println("No events received, skipping deletion check.")
		return
	}
	// Fetch all existing events from DB
	existingEvents, err := Events.GetAllEvents()
	if err != nil {
		log.Println("Error fetching events from DB:", err)
		return
	}

	// Create a map of received event UIDs
	receivedUIDs := make(map[string]bool)
	for _, event := range receivedEvents {
		receivedUIDs[event.UID] = true
	}

	// Check which events exist in DB but were NOT received
	for _, storedEvent := range existingEvents {
		if !receivedUIDs[storedEvent.UID] {
			log.Println("Event deleted:", storedEvent.UID)
			helpers.PrintAlert("event.deleted", storedEvent)
			// Delete the event from DB if needed
			err := Events.DeleteEvent(storedEvent.ID)
			if err != nil {
				log.Println("Error deleting event ", err)
				return
			}
		}
	}
}