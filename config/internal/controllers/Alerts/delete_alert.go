package Alerts

import (
	"database/sql"
	"github.com/google/uuid"
	"config/internal/services/Alerts"
	"net/http"
)

// DeleteAlertHandler deletes an alert by its ID.
// @Summary Delete an alert by its ID
// @Description This endpoint deletes an alert from the system using its unique identifier (ID).
// @Tags Alerts
// @Accept  json
// @Produce  json
// @Param alertId path string true "Alert ID"
// @Success 204 "No content"
// @Failure 400 "Invalid request"
// @Failure 404 "Alert not found"
// @Failure 500 "Internal server error"
// @Router /alerts/{alertId} [delete]
func DeleteAlertHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'ID de l'alerte depuis le contexte.
	alertId := r.Context().Value("alertId").(uuid.UUID)

	// Appeler le service pour supprimer l'alerte.
	err := Alerts.DeleteAlert(alertId)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Alert not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Répondre avec succès (aucune donnée dans la réponse).
	w.WriteHeader(http.StatusNoContent)
}