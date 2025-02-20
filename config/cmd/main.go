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
	"config/internal/controllers/AlertsControllers"
	"config/internal/controllers/RessourcesControllers"
	"config/internal/helpers"
	"github.com/go-chi/chi/v5"
)

func main() {
	// Initialisation de la base de données


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

	// Routes pour les ressources (resources)
	r.Route("/resources", func(r chi.Router) {
		r.Get("/", RessourcesControllers.GetAllResourcesHandler) // GET /config/resources
		r.Post("/", RessourcesControllers.CreateResourceHandler)  // POST /config/resources
		// Route pour les opérations sur une ressource spécifique
		r.Route("/{resourceId}", func(r chi.Router) { // Assurez-vous que le paramètre s'appelle "{id}"
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
			r.Route("/{alertId}", func(r chi.Router) {
				r.Use(controllers.Ctx("alertId")) // Utiliser controllers.Ctx ici
				r.Get("/", AlertsControllers.GetAlertByIDHandler)    // GET /config/alerts/{id}
				r.Put("/", AlertsControllers.UpdateAlertHandler)     // PUT /config/alerts/{id}
				r.Delete("/", AlertsControllers.DeleteAlertHandler)   // DELETE /config/alerts/{id}
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
