package Events

import (
	"encoding/json"
	"github.com/sirupsen/logrus"

	"timetable/internal/services/Events"
	"timetable/internal/models"
	"net/http"
)

// GetEventsHandler retrieves all events or events filtered by resource ID.
// @Summary Get all events or filter by resource ID
// @Description Fetches all events from the database, or filter events by the provided resource ID if specified.
// @Tags events
// @Accept json
// @Produce json
// @Param resourceId query string false "Resource ID" 
// @Success 200 {array} models.Event
// @Failure 400 {string} string "Invalid resource ID format"
// @Failure 500 {string} string "Internal server error"
// @Router /events [get]
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

