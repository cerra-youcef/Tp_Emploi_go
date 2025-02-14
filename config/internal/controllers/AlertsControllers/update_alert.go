package AlertsControllers

import (
"encoding/json"
"github.com/go-chi/chi/v5"
"github.com/sirupsen/logrus"
"cerra/tp_go/internal/models"
"cerra/tp_go/internal/services/AlertsSrv"
"github.com/google/uuid"


"net/http"
)

// UpdateAlertHandler gère la mise à jour d'une alerte existante.
func UpdateAlertHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	var updatedAlert models.Alert
	err = json.NewDecoder(r.Body).Decode(&updatedAlert)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = AlertsSrv.UpdateAlert(id, &updatedAlert)
	if err != nil {
		if customErr, ok := err.(*models.CustomError); ok {
			http.Error(w, customErr.Message, customErr.Code)
			return
		}
		logrus.Errorf("Error updating alert: %v", err)
		http.Error(w, "Failed to update alert", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
