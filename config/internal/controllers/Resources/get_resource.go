package Resources


import (
	"encoding/json"
"github.com/sirupsen/logrus"
"config/internal/models"
"config/internal/services/Resources"
"github.com/google/uuid"

"net/http"
)

// GetResourceByIDHandler retrieves a resource by its ID.
// @Summary Get a resource by its ID
// @Description This endpoint retrieves a resource by the provided resource ID.
// @Tags Resources
// @Param resourceId path string true "Resource ID" example("123e4567-e89b-12d3-a456-426614174000")
// @Success 200 {object} models.Resource "Resource found"
// @Failure 400 {string} string "Invalid resource ID"
// @Failure 404 {string} string "Resource not found"
// @Failure 500 {string} string "Internal server error"
// @Router /resources/{resourceId} [get]
func GetResourceByIDHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'ID de la ressource depuis le contexte.
	resourceId := r.Context().Value("resourceId").(uuid.UUID)

	// Appeler le service pour obtenir la ressource.
	resource, err := Resources.GetResourceByID(resourceId)
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
