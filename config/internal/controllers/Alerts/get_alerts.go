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
// @Param ucaID query string false "Filter alerts by UCA ID"
// @Success 200 {array} models.Alert "List of alerts"
// @Failure 500 {string} string "Internal server error"
// @Router /alerts [get]
func GetAllAlertsHandler(w http.ResponseWriter, r *http.Request) {
	ucaID := r.URL.Query().Get("ucaID")
	alerts, err := Alerts.GetAllAlerts(ucaID)
	if err != nil {
		logrus.Errorf("Error retrieving alerts: %v", err)
		http.Error(w, "Failed to retrieve alerts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(alerts)
}
