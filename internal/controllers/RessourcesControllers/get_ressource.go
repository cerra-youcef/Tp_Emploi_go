package RessourcesControllers


import (
	"encoding/json"
"github.com/sirupsen/logrus"
"middleware/example/internal/models"
"middleware/example/internal/services/RessourcesSrv"
"github.com/google/uuid"

"net/http"
)


// GetResourceByIDHandler récupère une ressource par son ID.
func GetResourceByIDHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'ID de la ressource depuis le contexte.
	resourceId := r.Context().Value("resourceId").(uuid.UUID)

	// Appeler le service pour obtenir la ressource.
	resource, err := RessourcesSrv.GetResourceByID(resourceId)
	if err != nil {
		if customErr, ok := err.(*models.CustomError); ok {
			http.Error(w, customErr.Message, customErr.Code)
			return
		}
		logrus.Errorf("Error retrieving resource: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Vérifier si la ressource existe.
	if resource == nil {
		http.Error(w, "Resource not found", http.StatusNotFound)
		return
	}

	// Répondre avec la ressource au format JSON.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resource)
}
