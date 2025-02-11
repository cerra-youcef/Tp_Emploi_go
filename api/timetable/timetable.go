package timetable

import (
	"net/http"

	"middleware/example/internal/controllers" // Import du package controllers
	"middleware/example/internal/controllers/EventsControllers"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes() http.Handler {
	r := chi.NewRouter()

	// Routes pour les événements (events)
	r.Route("/events", func(r chi.Router) {
		r.Get("/", EventsControllers.GetAllEventsHandler) // GET /timetable/events
		r.Post("/", EventsControllers.CreateEventHandler) // POST /timetable/events
		r.Route("/{id}", func(r chi.Router) {
			r.Use(controllers.Ctx("eventId")) // Utiliser controllers.Ctx ici
			r.Get("/", EventsControllers.GetEventByIDHandler)    // GET /timetable/events/{id}
			r.Put("/", EventsControllers.UpdateEventHandler)     // PUT /timetable/events/{id}
			r.Delete("/", EventsControllers.DeleteEventHandler)   // DELETE /timetable/events/{id}
		})
	})

	// Route pour récupérer les événements par resource ID
	r.Get("/events/resources/{resourceId}", EventsControllers.GetEventsByResourceIDHandler) // GET /timetable/events/resources/{resourceId}

	return r
}