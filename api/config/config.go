package timetable

import (
	"net/http"

	"middleware/example/internal/controllers/EventsControllers"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(r chi.Router) {
	// Routes pour les événements (events)
	r.Route("/events", func(r chi.Router) {
		r.Get("/", EventsControllers.GetAllEventsHandler) // GET /timetable/events
		r.Post("/", EventsControllers.CreateEventHandler) // POST /timetable/events
		r.Route("/{id}", func(r chi.Router) {
			r.Use(EventsControllers.Ctx("eventId")) // Middleware pour injecter l'ID dans le contexte
			r.Get("/", EventsControllers.GetEventByIDHandler)    // GET /timetable/events/{id}
			r.Put("/", EventsControllers.UpdateEventHandler)     // PUT /timetable/events/{id}
			r.Delete("/", EventsControllers.DeleteEventHandler)   // DELETE /timetable/events/{id}
		})
	})

	// Route pour récupérer les événements par resource ID
	r.Get("/events/resources/{resourceId}", EventsControllers.GetEventsByResourceIDHandler) // GET /timetable/events/resources/{resourceId}
}