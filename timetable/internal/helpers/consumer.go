package helpers

import (
	"log"
	"timetable/internal/models"
)

func PrintAlert(alertType string, event models.Event, changes ...map[string]string){
	natsMessage := map[string]interface{}{
		"type":  alertType,
		"event": event,
	}

	if len(changes) > 0 {
		natsMessage["changes"] = changes[0]
	}

	//message, _ := json.Marshal(natsMessage)
	log.Println(natsMessage["type"])
	log.Println(natsMessage["changes"])
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
