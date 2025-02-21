package Alerts

import (
	"encoding/json"
	"github.com/sirupsen/logrus"

	"config/internal/services/Alerts"
	"github.com/google/uuid"

	"net/http"
)

// GetAlertByIDHandler retrieves an alert by its ID.
// @Summary Get an alert by its ID
// @Description This endpoint retrieves an alert by its unique identifier (ID).
// @Tags Alerts
// @Accept  json
// @Produce  json
// @Param alertId path string true "Alert ID"
// @Success 200 {object} models.Alert "Alert found"
// @Failure 400 {string} string "Invalid request"
// @Failure 404 {string} string "Alert not found"
// @Failure 500 {string} string "Internal server error"
// @Router /alerts/{alertId} [get]
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

