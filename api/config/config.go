package config

import (
	"net/http"

	"middleware/example/internal/controllers" // Import du package controllers
	"middleware/example/internal/controllers/AlertsControllers"
	"middleware/example/internal/controllers/RessourcesControllers"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes() http.Handler {
	r := chi.NewRouter()

	// Routes pour les ressources (resources)
	r.Route("/resources", func(r chi.Router) {
		r.Get("/", RessourcesControllers.GetAllResourcesHandler) // GET /config/resources
		r.Post("/", RessourcesControllers.CreateResourceHandler)  // POST /config/resources
		// Route pour les opérations sur une ressource spécifique
		r.Route("/{id}", func(r chi.Router) { // Assurez-vous que le paramètre s'appelle "{id}"
			r.Use(controllers.Ctx("resourceId")) // Utiliser controllers.Ctx ici
			r.Get("/", RessourcesControllers.GetResourceByIDHandler)    // GET /config/resources/{id}
			r.Put("/", RessourcesControllers.UpdateResourceHandler)     // PUT /config/resources/{id}
			r.Delete("/", RessourcesControllers.DeleteResourceHandler)   // DELETE /config/resources/{id}
		})
	})

	// Routes pour les alertes (alerts)
	r.Route("/alerts", func(r chi.Router) {
		r.Get("/", AlertsControllers.GetAllAlertsHandler) // GET /config/alerts
		r.Post("/", AlertsControllers.CreateAlertHandler)  // POST /config/alerts
		r.Route("/{id}", func(r chi.Router) {
			r.Use(controllers.Ctx("alertId")) // Utiliser controllers.Ctx ici
			r.Get("/", AlertsControllers.GetAlertByIDHandler)    // GET /config/alerts/{id}
			r.Put("/", AlertsControllers.UpdateAlertHandler)     // PUT /config/alerts/{id}
			r.Delete("/", AlertsControllers.DeleteAlertHandler)   // DELETE /config/alerts/{id}
		})
	})

	return r
}