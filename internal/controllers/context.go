package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	"middleware/example/internal/helpers"
	"middleware/example/internal/models"
	"net/http"
)

// Ctx injecte un ID UUID et la connexion DB dans le contexte.
func Ctx(idParam string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Récupérer l'ID depuis les paramètres de l'URL.
			idStr := chi.URLParam(r, idParam)
			if idStr == "" {
				logrus.Errorf("Missing %s in URL", idParam)
				errorResponse(w, fmt.Sprintf("Invalid request: missing %s", idParam), http.StatusBadRequest)
				return
			}

			// Valider et convertir l'ID en UUID.
			id, err := uuid.FromString(idStr)
			if err != nil {
				logrus.Errorf("Invalid %s: %s", idParam, idStr)
				customError := &models.CustomError{
					Message: fmt.Sprintf("Cannot parse %s (%s) as UUID", idParam, idStr),
					Code:    http.StatusUnprocessableEntity,
				}
				errorResponse(w, customError.Message, customError.Code)
				return
			}

			// Ouvrir une connexion à la base de données.
			db, err := helpers.OpenDB()
			if err != nil {
				logrus.Errorf("Failed to open database connection: %v", err)
				errorResponse(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			// Ajouter l'ID et la connexion DB au contexte.
			ctx := context.WithValue(r.Context(), idParam, id)
			ctx = context.WithValue(ctx, "db", db)

			// Passer au handler suivant.
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// errorResponse simplifie l'envoi d'une réponse JSON en cas d'erreur.
func errorResponse(w http.ResponseWriter, message string, code int) {
	w.WriteHeader(code)
	body, _ := json.Marshal(map[string]string{"error": message})
	_, _ = w.Write(body)
}