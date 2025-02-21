package Events

import (
	"github.com/google/uuid"

	"encoding/json"
	"github.com/sirupsen/logrus"
	"timetable/internal/models"
	"timetable/internal/services/Events"
	"net/http"
)

// GetEventByIDHandler retrieves an event by its ID.
// @Summary Get event by ID
// @Description Retrieves an event by its unique ID from the database
// @Tags events
// @Accept json
// @Produce json
// @Param eventId path string true "Event ID"
// @Success 200 {object} models.Event
// @Failure 400 {string} string "Invalid Event ID"
// @Failure 404 {string} string "Event not found"
// @Failure 500 {string} string "Internal server error"
// @Router /events/{eventId} [get]
func GetEventByIDHandler(w http.ResponseWriter, r *http.Request) {
	eventID := r.Context().Value("eventId").(uuid.UUID)

	event, err := Events.GetEventByID(eventID)
	if err != nil {
		if customErr, ok := err.(*models.CustomError); ok {
			http.Error(w, customErr.Message, customErr.Code)
			return
		}
		logrus.Errorf("Error retrieving event: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if event == nil {
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(event)
}

