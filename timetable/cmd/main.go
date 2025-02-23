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
	"github.com/swaggo/http-swagger"
	"github.com/go-chi/cors"
	_ "timetable/api"
	"timetable/internal/nats"
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

	//Store DB in context
	ctx := context.WithValue(context.Background(), "db", db)

	//Start NATS Consumer in a Goroutine with context
	go func() {
		log.Println("Starting NATS Consumer...")
		js, nc, err := natsConsumer.ConnectToNATS()
		if err != nil {
			log.Printf("Error connecting to NATS: %v", err)
			return //to keep the api runnig
		}
		defer nc.Close()

		consumer, err := natsConsumer.EventConsumer(js)
		if err != nil {
			log.Printf("Error creating NATS Consumer: %v", err)
			return //to keep the api runnig
		}

		err = natsConsumer.Consume(ctx, *consumer) // ✅ Pass the context with DB
		if err != nil {
			log.Printf("Error consuming messages: %v", err)
			return //to keep the api runnig
		}
	}()


	// Création du routeur Chi
	r := chi.NewRouter()

	r.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Allow all origins
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},                    // Allow specific HTTP methods
		AllowedHeaders:   []string{"Content-Type", "Authorization"},                   // Allow specific headers
		AllowCredentials: true,                                                      // Allow credentials (cookies, authorization headers)
	}).Handler)

	// Routes pour les événements (events)
	r.Route("/events", func(r chi.Router) {
		r.Get("/", Events.GetEventsHandler) // GET /timetable/events
		r.Post("/", Events.CreateEventHandler) // POST /timetable/events
		r.Route("/{eventId}", func(r chi.Router) {
			r.Use(controllers.Ctx("eventId")) // Utiliser controllers.Ctx ici
			r.Get("/", Events.GetEventByIDHandler)    // GET /timetable/events/{id}
		})
	})

	// Swagger UI Route
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), // Point to your swagger.json
	))

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