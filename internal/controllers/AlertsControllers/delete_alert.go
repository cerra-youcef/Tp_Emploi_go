package AlertsControllers

import (
	"database/sql"

	"github.com/google/uuid"



	"middleware/example/internal/services/AlertsSrv"


	"net/http"
)

// DeleteAlertHandler supprime une alerte par son ID.
func DeleteAlertHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'ID de l'alerte depuis le contexte.
	alertId := r.Context().Value("alertId").(uuid.UUID)

	// Appeler le service pour supprimer l'alerte.
	err := AlertsSrv.DeleteAlert(alertId)
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