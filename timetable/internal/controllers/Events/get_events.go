package Events

import (
	"encoding/json"
	"github.com/sirupsen/logrus"

	"timetable/internal/services/Events"
	"timetable/internal/models"
	"net/http"
)

// GetAllEventsHandler récupère tous les événements.
func GetEventsHandler(w http.ResponseWriter, r *http.Request) {
	resourceId := r.URL.Query().Get("resourceId")
	var events []models.Event
	var err error
	if resourceId != "" {
		// Fetch events filtered by resourceId
		events, err = Events.GetEventsByResourceID(resourceId)
	} else {
		// Fetch all events
		events, err = Events.GetAllEvents()
	}

	if err != nil {
		logrus.Errorf("Error retrieving events: %v", err)
		http.Error(w, "Failed to retrieve events", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

