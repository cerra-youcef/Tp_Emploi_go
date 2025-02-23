package helpers

import (
	"log"
	"timetable/internal/models"
	"github.com/google/uuid"

)

func PrintAlert(alertType string, event models.Event, changes ...map[string]string){

	alert := models.Alert{
		ID:    uuid.New(),  // Generate a new UUID for the alert
		Type:  alertType,
		Event: event.ID,    // Assuming your Event model has an ID field, otherwise update this
	}

	if len(changes) > 0 {
		alert.Changes = changes[0]
	}

	//message, _ := json.Marshal(natsMessage)
	log.Printf("Alert Type: %s\n", alert.Type)
	log.Printf("Event ID: %s\n", alert.Event)
	log.Printf("Changes: %v\n", alert.Changes)
}



func DetectChanges(oldEvent, newEvent models.Event) map[string]string {
	changes := make(map[string]string)

	if oldEvent.Name != newEvent.Name {
		changes["name"] = newEvent.Name
	}
	if oldEvent.Description != newEvent.Description {
		changes["description"] = newEvent.Description
	}
	if oldEvent.Location != newEvent.Location {
		changes["location"] = newEvent.Location
	}
	if oldEvent.Start != newEvent.Start {
		changes["start"] = newEvent.Start
	}
	if oldEvent.End != newEvent.End {
		changes["end"] = newEvent.End
	}

	return changes
}
