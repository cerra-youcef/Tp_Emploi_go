package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"timetable/internal/helpers"
	"timetable/internal/controllers/Events"
	"timetable/internal/controllers"
	"github.com/go-chi/chi/v5"
)

func main() {
	//Initialisation de la base de données
	db, err := helpers.OpenDB()
	if err != nil {
		log.Fatalf("Error while opening database: %s", err.Error())
	}
	defer helpers.CloseDB(db)
	
	// Vérification initiale des tables dans la base de données
	if err := helpers.InitializeDB(db); err != nil {
		log.Fatalf("Error initializing database: %s", err.Error())
	}
	



	// Création du routeur Chi
	r := chi.NewRouter()

	// Routes pour les événements (events)
	r.Route("/events", func(r chi.Router) {
		r.Get("/", Events.GetAllEventsHandler) // GET /timetable/events
		r.Post("/", Events.CreateEventHandler) // POST /timetable/events
		r.Route("/{eventId}", func(r chi.Router) {
			r.Use(controllers.Ctx("eventId")) // Utiliser controllers.Ctx ici
			r.Get("/", Events.GetEventByIDHandler)    // GET /timetable/events/{id}
		})
	})

	// Passage de la connexion DB au contexte pour une utilisation globale
	ctx := context.WithValue(context.Background(), "db", db)

	// Démarrage du serveur HTTP avec un gestionnaire de contexte
	server := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Canal pour écouter les signaux système
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// Démarrer le serveur dans un goroutine
	go func() {
		log.Println("Web server started. Listening on :8080")
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Attendre un signal dearrêt
	sig := <-shutdown
	log.Printf("Received shutdown signal: %v", sig)

	// Arrêter proprement le serveur
	log.Println("Shutting down server...")
	ctxShutdown, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctxShutdown); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server stopped gracefully.")
}