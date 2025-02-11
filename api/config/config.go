package config

import (
	"net/http"

	"middleware/example/internal/controllers/AlertsControllers"
	"middleware/example/internal/controllers/RessourcesControllers"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(r chi.Router) {
	// Routes pour les ressources (resources)
	r.Route("/resources", func(r chi.Router) {
		r.Get("/", RessourcesControllers.GetAllResourcesHandler) // GET /config/resources
		r.Post("/", RessourcesControllers.CreateResourceHandler)  // POST /config/resources
		r.Route("/{id}", func(r chi.Router) {
			r.Use(RessourcesControllers.Ctx("resourceId")) // Middleware pour injecter l'ID dans le contexte
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
			r.Use(AlertsControllers.Ctx("alertId")) // Middleware pour injecter l'ID dans le contexte
			r.Get("/", AlertsControllers.GetAlertByIDHandler)    // GET /config/alerts/{id}
			r.Put("/", AlertsControllers.UpdateAlertHandler)     // PUT /config/alerts/{id}
			r.Delete("/", AlertsControllers.DeleteAlertHandler)   // DELETE /config/alerts/{id}
		})
	})
}