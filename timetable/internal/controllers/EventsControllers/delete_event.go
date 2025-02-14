package EventsControllers

import (
	"github.com/google/uuid"


	"github.com/sirupsen/logrus"
	"cerra/tp_go/internal/services/EventsSrv"
	"net/http"
	"database/sql"

)


// DeleteEventHandler supprime un événement existant.
func DeleteEventHandler(w http.ResponseWriter, r *http.Request) {
	eventID := r.Context().Value("eventId").(uuid.UUID)

	err := EventsSrv.DeleteEvent(eventID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Event not found", http.StatusNotFound)
			return
		}
		logrus.Errorf("Error deleting event: %v", err)
		http.Error(w, "Failed to delete event", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
