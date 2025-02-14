package EventsControllers

import (
	"github.com/google/uuid"

	"encoding/json"
	"github.com/sirupsen/logrus"
	"cerra/tp_go/internal/models"
	"cerra/tp_go/internal/services/EventsSrv"
	"net/http"

)

// GetEventByIDHandler récupère un événement par son ID.
func GetEventByIDHandler(w http.ResponseWriter, r *http.Request) {
	eventID := r.Context().Value("eventId").(uuid.UUID)

	event, err := EventsSrv.GetEventByID(eventID)
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

