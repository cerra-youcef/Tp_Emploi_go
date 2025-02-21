package Events

import (
	"github.com/google/uuid"

	"encoding/json"
	"github.com/sirupsen/logrus"
	"timetable/internal/models"
	"timetable/internal/services/Events"
	"net/http"
)

// CreateEventHandler creates a new event.
// @Summary Create a new event
// @Description Creates a new event in the system. The event must include details like the event name, description, start time, and resource associations.
// @Tags events
// @Accept json
// @Produce json
// @Param event body models.Event true "Event to create"
// @Success 201 {object} models.Event "Event created successfully"
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Internal server error"
// @Router /events [post]
func CreateEventHandler(w http.ResponseWriter, r *http.Request) {
	var newEvent models.Event
	err := json.NewDecoder(r.Body).Decode(&newEvent)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if newEvent.ID == uuid.Nil {
		newEvent.ID = uuid.New()
	}

	err = Events.CreateEvent(&newEvent)
	if err != nil {
		if customErr, ok := err.(*models.CustomError); ok {
			http.Error(w, customErr.Message, customErr.Code)
			return
		}
		logrus.Errorf("Error creating event: %v", err)
		http.Error(w, "Failed to create event", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newEvent)
}

