package Resources

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"config/internal/models"
	"config/internal/services/Resources"
	"net/http"
	"github.com/google/uuid"

)

// CreateResourceHandler creates a new resource.
// @Summary Create a new resource
// @Description This endpoint creates a new resource in the system.
// @Tags Resources
// @Accept json
// @Produce json
// @Param resource body models.Resource true "Resource object"
// @Success 201 {object} models.Resource "The newly created resource"
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Failed to create resource"
// @Router /resources [post]
func CreateResourceHandler(w http.ResponseWriter, r *http.Request) {
	var newResource models.Resource
	err := json.NewDecoder(r.Body).Decode(&newResource)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if newResource.ID == uuid.Nil {
		newResource.ID = uuid.New()
	}

	err = Resources.CreateResource(&newResource)
	if err != nil {
		if customErr, ok := err.(*models.CustomError); ok {
			http.Error(w, customErr.Message, customErr.Code)
			return
		}
		logrus.Errorf("Error creating resource: %v", err)
		http.Error(w, "Failed to create resource", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newResource)
}

