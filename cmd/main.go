package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"middleware/example/internal/controllers"
	"middleware/example/internal/controllers/AlertsControllers"
	"middleware/example/internal/controllers/EventsControllers"
	"middleware/example/internal/controllers/RessourcesControllers"
	"middleware/example/internal/helpers"
)

func main() {
	r := chi.NewRouter()

	// Routes pour les Ressources
	r.Route("/resources", func(r chi.Router) {
		r.Get("/", RessourcesControllers.GetAllResourcesHandler) // GET /resources
		r.Post("/", RessourcesControllers.CreateResourceHandler)  // POST /resources

		r.Route("/{id}", func(r chi.Router) {
			r.Use(controllers.Ctx("resourceId")) // Middleware pour injecter l'ID dans le contexte
			r.Get("/", RessourcesControllers.GetResourceByIDHandler)    // GET /resources/{id}
			r.Put("/", RessourcesControllers.UpdateResourceHandler)     // PUT /resources/{id}
			r.Delete("/", RessourcesControllers.DeleteResourceHandler)   // DELETE /resources/{id}
		})
	})

	// Routes pour les Alertes
	r.Route("/alerts", func(r chi.Router) {
		r.Get("/", AlertsControllers.GetAllAlertsHandler) // GET /alerts
		r.Post("/", AlertsControllers.CreateAlertHandler)  // POST /alerts

		r.Route("/{id}", func(r chi.Router) {
			r.Use(controllers.Ctx("alertId")) // Middleware pour injecter l'ID dans le contexte
			r.Get("/", AlertsControllers.GetAlertByIDHandler)    // GET /alerts/{id}
			r.Put("/", AlertsControllers.UpdateAlertHandler)     // PUT /alerts/{id}
			r.Delete("/", AlertsControllers.DeleteAlertHandler)   // DELETE /alerts/{id}
		})
	})

	// Routes pour les Événements
	r.Route("/events", func(r chi.Router) {

		r.Get("/", EventsControllers.GetAllEventsHandler) // GET /events
		r.Post("/", EventsControllers.CreateEventHandler) // POST /events

		r.Route("/{id}", func(r chi.Router) {
			r.Use(controllers.Ctx("eventId")) // Middleware pour injecter l'ID dans le contexte
			r.Get("/", EventsControllers.GetEventByIDHandler)    // GET /events/{id}
			r.Put("/", EventsControllers.UpdateEventHandler)     // PUT /events/{id}
			r.Delete("/", EventsControllers.DeleteEventHandler)   // DELETE /events/{id}
		})
	})

	// Initialisation de la base de données
	db, err := helpers.OpenDB()
	if err != nil {
		logrus.Fatalf("Error while opening database: %s", err.Error())
	}
	defer helpers.CloseDB(db)

	// Initialiser les tables (si elles n'existent pas déjà)
	if err := helpers.InitializeDB(db); err != nil {
		logrus.Fatalf("Error initializing database: %s", err.Error())
	}

	// Démarrage du serveur HTTP
	logrus.Info("[INFO] Web server started. Now listening on *:8080")
	logrus.Fatalln(http.ListenAndServe(":8080", r))
}