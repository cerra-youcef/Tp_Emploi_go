package Alerts

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"

	"config/internal/services/Alerts"
	"github.com/google/uuid"


	"net/http"
)

// GetAlertByIDHandler gère la récupération d'une alerte par son ID.
func GetAlertByIDHandler(w http.ResponseWriter, r *http.Request) {
idParam := chi.URLParam(r, "id")
id, err := uuid.Parse(idParam)
if err != nil {
http.Error(w, "Invalid UUID", http.StatusBadRequest)
return
}

alert, err := Alerts.GetAlertByID(id)
if err != nil {
logrus.Errorf("Error retrieving alert: %v", err)
http.Error(w, "Internal server error", http.StatusInternalServerError)
return
}

if alert == nil {
http.Error(w, "Alert not found", http.StatusNotFound)
return
}

w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(alert)
}

