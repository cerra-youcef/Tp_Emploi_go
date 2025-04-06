package helpers

import (
	//	"log"
	"time"
	"timetable/internal/models"

	"github.com/google/uuid"
)

func CreateAlert(alertType string, event models.Event, changes ...map[string]string) models.Alert {

	t, _ := time.Parse("20060102T150405Z", event.Start)
	formatted_date := t.Format("2006-01-02 15:04:05")

	alert := models.Alert{
		ID:        uuid.New(),
		Type:      alertType,
		Event:     event.Name,
		Resources: event.Resources,
		Location:  event.Location,
		Start:     formatted_date,
	}

	if len(changes) > 0 {
		alert.Changes = changes[0]
	}
	return alert
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
