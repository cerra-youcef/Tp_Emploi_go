package Alerts

import (
	"encoding/json"
	"github.com/sirupsen/logrus"

	"config/internal/services/Alerts"
	"github.com/google/uuid"

	"net/http"
)

// GetAlertByIDHandler gère la récupération d'une alerte par son ID.
func GetAlertByIDHandler(w http.ResponseWriter, r *http.Request) {

	alertIdRaw := r.Context().Value("alertId") 
	if alertIdRaw == nil {
		logrus.Errorf("Alert ID not found in context")
		http.Error(w, "Missing alert in context", http.StatusBadRequest)
		return
	}

	alertId, ok := alertIdRaw.(uuid.UUID)
	if !ok {
		logrus.Errorf("Alert ID is not of type uuid.UUID")
		http.Error(w, "Invalid alert ID type", http.StatusBadRequest)
		return
	}

	alert, err := Alerts.GetAlertByID(alertId)
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

