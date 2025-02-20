package EventsControllers

import (
"encoding/json"
"github.com/go-chi/chi/v5"

"net/http"
"database/sql"

"timetable/internal/services/EventsSrv"
)

// GetEventsByResourceIDHandler récupère les événements associés à une ressource spécifique.
func GetEventsByResourceIDHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer la connexion à la base de données depuis le contexte
	db := r.Context().Value("db").(*sql.DB)

	// Récupérer l'ID de la ressource depuis les paramètres de l'URL
	resourceID := chi.URLParam(r, "resourceId")

	// Appeler le service pour obtenir les événements
	events, err := EventsSrv.GetEventsByResourceID(db, resourceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Répondre avec les événements au format JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}
