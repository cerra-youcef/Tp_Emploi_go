package helpers

import (
	"database/sql"
	"context"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func InitializeDB(db *sql.DB) error {
	// Script SQL pour créer la table `collections` (ressources).
	resourcesTableQuery := `
        CREATE TABLE IF NOT EXISTS resources (
            id TEXT PRIMARY KEY,
            uca_id INTEGER NOT NULL,
            name TEXT NOT NULL
        );
    `
	// Script SQL pour créer la table `alerts`.
	alertsTableQuery := `
        CREATE TABLE IF NOT EXISTS alerts (
            id TEXT PRIMARY KEY,
            email TEXT NOT NULL,
            is_all BOOLEAN NOT NULL,
            resource_id TEXT
        );
    `
	// Script SQL pour créer la table `events`.
	eventsTableQuery := `
        CREATE TABLE IF NOT EXISTS events (
            id TEXT PRIMARY KEY,
            resource_ids TEXT NOT NULL,
            uid TEXT NOT NULL,
            name TEXT NOT NULL,
            start TEXT NOT NULL
        );
    `
	// Exécuter les requêtes SQL pour créer les tables.
	queries := []string{resourcesTableQuery, alertsTableQuery, eventsTableQuery}
	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetDBFromContext récupère la connexion DB depuis le contexte.
func GetDBFromContext(ctx context.Context) *sql.DB {
	db, ok := ctx.Value("db").(*sql.DB)
	if !ok {
		panic("Database connection not found in context")
	}
	return db
}

// OpenDB ouvre une connexion à la base de données SQLite.
func OpenDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "file:timetable.db?cache=shared&mode=rwc")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
		return nil, err
	}

	// Limiter le nombre de connexions ouvertes (utile pour SQLite).
	db.SetMaxOpenConns(1)

	// Initialiser la base de données (créer les tables si elles n'existent pas).
	if err := InitializeDB(db); err != nil {
		log.Fatalf("Error initializing database: %v", err)
		return nil, err
	}

	return db, nil
}

// CloseDB ferme la connexion à la base de données.
func CloseDB(db *sql.DB) {
	if err := db.Close(); err != nil {
		log.Printf("Error closing database: %v", err)
	}
}

