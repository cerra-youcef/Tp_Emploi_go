package helpers

import (
	"database/sql"
	"context"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func InitializeDB(db *sql.DB) error {

	// Script SQL pour créer la table `events`.
	eventsTableQuery := `
        CREATE TABLE IF NOT EXISTS events (
            id TEXT PRIMARY KEY,
            resources TEXT NOT NULL,
            uid TEXT NOT NULL,
            name TEXT NOT NULL,
			description TEXT,
            start TEXT NOT NULL,
			end TEXT,
			location TEXT,
			UpdatedAt TEXT
        );
    `

	_, err := db.Exec(eventsTableQuery)
	if err != nil {
		return err
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

