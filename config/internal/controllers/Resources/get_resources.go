package Resources

import (

"encoding/json"
"github.com/sirupsen/logrus"
"config/internal/services/Resources"
"net/http"
)

// GetAllResourcesHandler retrieves all resources.
// @Summary Get all resources
// @Description This endpoint retrieves a list of all available resources.
// @Tags Resources
// @Success 200 {array} models.Resource "List of resources"
// @Failure 500 {string} string "Failed to retrieve resources"
// @Router /resources [get]
func GetAllResourcesHandler(w http.ResponseWriter, r *http.Request) {
	resources, err := Resources.GetAllResources()
	if err != nil {
		logrus.Errorf("Error retrieving resources: %v", err)
		http.Error(w, "Failed to retrieve resources", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resources)
}
