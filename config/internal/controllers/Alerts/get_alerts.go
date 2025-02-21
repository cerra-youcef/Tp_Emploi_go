package Alerts

import (
"encoding/json"
"github.com/sirupsen/logrus"

"config/internal/services/Alerts"

"net/http"
)

// GetAllAlertsHandler gère la récupération de toutes les alertes.
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
