package main

import (
	"log"
	"net/http"

	"middleware/example/api/config"
	"middleware/example/api/timetable"
	"middleware/example/internal/helpers"
	"github.com/go-chi/chi/v5"
)

func main() {
	// Initialisation de la base de données
	db, err := helpers.OpenDB()
	if err != nil {
		log.Fatalf("Error while opening database: %s", err.Error())
	}
	defer helpers.CloseDB(db)

	// Initialiser les tables (si elles n'existent pas déjà)
	if err := helpers.InitializeDB(db); err != nil {
		log.Fatalf("Error initializing database: %s", err.Error())
	}

	// Création du routeur Chi
	r := chi.NewRouter()

	// Ajout des routes pour l'API "Config"
	r.Mount("/config", config.SetupRoutes())

	// Ajout des routes pour l'API "Timetable"
	r.Mount("/timetable", timetable.SetupRoutes())

	// Démarrage du serveur HTTP
	log.Println("Web server started. Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}