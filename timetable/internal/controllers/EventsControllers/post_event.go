package EventsControllers

import (
	"github.com/google/uuid"

	"encoding/json"
	"github.com/sirupsen/logrus"
	"cerra/tp_go/internal/models"
	"cerra/tp_go/internal/services/EventsSrv"
	"net/http"
)

// CreateEventHandler crée un nouvel événement.
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

	err = EventsSrv.CreateEvent(&newEvent)
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

