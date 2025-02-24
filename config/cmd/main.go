package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"config/internal/controllers"
	"config/internal/controllers/Alerts"
	"config/internal/controllers/Resources"
	"config/internal/helpers"
	"github.com/go-chi/chi/v5"
	"github.com/swaggo/http-swagger"
	"github.com/go-chi/cors"
	_ "config/api"
	"config/internal/nats"
)

func main() {
	db, err := helpers.OpenDB()
	if err != nil {
		log.Fatalf("Error while opening database: %s", err.Error())
	}
	defer helpers.CloseDB(db)

	// Vérification initiale des tables dans la base de données
	if err := helpers.InitializeDB(db); err != nil {
		log.Fatalf("Error initializing database: %s", err.Error())
	}

	//Start NATS Consumer in a Goroutine with context
	go func() {
		log.Println("Starting NATS Consumer...")
		js, nc, err := natsConsumer.ConnectToNATS()
		if err != nil {
			log.Printf("Error connecting to NATS: %v", err)
			return //to keep the api runnig
		}
		defer nc.Close()

		consumer, err := natsConsumer.AlertConsumer(js)
		if err != nil {
			log.Printf("Error creating NATS Consumer: %v", err)
			return //to keep the api runnig
		}

		err = natsConsumer.Consume(*consumer) // Pass the context with DB
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

	// Routes pour les ressources (resources)
	r.Route("/resources", func(r chi.Router) {
		r.Get("/", Resources.GetAllResourcesHandler) // GET /config/resources
		r.Post("/", Resources.CreateResourceHandler)  // POST /config/resources
		// Route pour les opérations sur une ressource spécifique
		r.Route("/{resourceId}", func(r chi.Router) { // Assurez-vous que le paramètre s'appelle "{id}"
			r.Use(controllers.Ctx("resourceId")) // Utiliser controllers.Ctx ici
			r.Get("/", Resources.GetResourceByIDHandler)    // GET /config/resources/{id}
			r.Put("/", Resources.UpdateResourceHandler)     // PUT /config/resources/{id}
			r.Delete("/", Resources.DeleteResourceHandler)   // DELETE /config/resources/{id}
		})
	})

	// Routes pour les alertes (alerts)
	r.Route("/alerts", func(r chi.Router) {
		r.Get("/", Alerts.GetAllAlertsHandler) // GET /config/alerts
		r.Post("/", Alerts.CreateAlertHandler)  // POST /config/alerts
		r.Route("/{alertId}", func(r chi.Router) {
			r.Use(controllers.Ctx("alertId")) // Utiliser controllers.Ctx ici
			r.Get("/", Alerts.GetAlertByIDHandler)    // GET /config/alerts/{id}
			r.Put("/", Alerts.UpdateAlertHandler)     // PUT /config/alerts/{id}
			r.Delete("/", Alerts.DeleteAlertHandler)   // DELETE /config/alerts/{id}
		})
	})

	// Swagger UI Route
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), // Point to your swagger.json
	))

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
