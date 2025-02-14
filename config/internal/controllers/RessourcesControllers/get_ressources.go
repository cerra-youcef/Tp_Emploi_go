package RessourcesControllers

import (

"encoding/json"
"github.com/sirupsen/logrus"
"cerra/tp_go/internal/services/RessourcesSrv"
"net/http"
)

// GetAllResourcesHandler récupère toutes les ressources.
func GetAllResourcesHandler(w http.ResponseWriter, r *http.Request) {
	resources, err := RessourcesSrv.GetAllResources()
	if err != nil {
		logrus.Errorf("Error retrieving resources: %v", err)
		http.Error(w, "Failed to retrieve resources", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resources)
}
