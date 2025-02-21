package Events

import (
"encoding/json"
"github.com/sirupsen/logrus"
"timetable/internal/models"
"timetable/internal/services/Events"
"net/http"
"github.com/google/uuid"

)

// UpdateEventHandler met à jour un événement existant.
func UpdateEventHandler(w http.ResponseWriter, r *http.Request) {
	eventID := r.Context().Value("eventId").(uuid.UUID)

	var updatedEvent models.Event
	err := json.NewDecoder(r.Body).Decode(&updatedEvent)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if updatedEvent.ID != eventID {
		http.Error(w, "ID mismatch between URL and request body", http.StatusBadRequest)
		return
	}

	err = Events.UpdateEvent(eventID, &updatedEvent)
	if err != nil {
		if customErr, ok := err.(*models.CustomError); ok {
			http.Error(w, customErr.Message, customErr.Code)
			return
		}
		logrus.Errorf("Error updating event: %v", err)
		http.Error(w, "Failed to update event", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}


