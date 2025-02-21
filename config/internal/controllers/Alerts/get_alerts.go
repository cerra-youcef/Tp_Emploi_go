package Alerts

import (
"encoding/json"
"github.com/sirupsen/logrus"

"config/internal/services/Alerts"

"net/http"
)

// GetAllAlertsHandler retrieves all alerts.
// @Summary Get all alerts
// @Description This endpoint retrieves a list of all alerts.
// @Tags Alerts
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Alert "List of alerts"
// @Failure 500 {string} string "Internal server error"
// @Router /alerts [get]
func GetAllAlertsHandler(w http.ResponseWriter, r *http.Request) {
	alerts, err := Alerts.GetAllAlerts()
	if err != nil {
		logrus.Errorf("Error retrieving alerts: %v", err)
		http.Error(w, "Failed to retrieve alerts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(alerts)
}
