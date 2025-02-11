package EventsControllers

import (
	"encoding/json"
	"github.com/sirupsen/logrus"

	"middleware/example/internal/services/EventsSrv"
	"net/http"

)

// GetAllEventsHandler récupère tous les événements.
func GetAllEventsHandler(w http.ResponseWriter, r *http.Request) {
	events, err := EventsSrv.GetAllEvents()
	if err != nil {
		logrus.Errorf("Error retrieving events: %v", err)
		http.Error(w, "Failed to retrieve events", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

