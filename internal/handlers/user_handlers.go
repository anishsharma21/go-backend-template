package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/anishsharma21/go-web-dev-template/internal/queries"
	"github.com/anishsharma21/go-web-dev-template/internal/types/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func AddNewUser(dbPool *pgxpool.Pool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ctx = r.Context()
		var user models.User

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			slog.Error("Failed to decode request body", "error", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		query := "INSERT INTO users (clerk_id) VALUES ($1)"
		result, err := dbPool.Exec(ctx, query, user.ClerkID)
		if err != nil {
			slog.Error("Failed to insert user", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if result.RowsAffected() == 0 {
			slog.Error("No rows affected")
			http.Error(w, "User not created", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	})
}

func GetUsers(dbPool *pgxpool.Pool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ctx = r.Context()
		query := "SELECT id, clerk_id, created_at FROM users"

		rows, err := dbPool.Query(ctx, query)
		if err != nil {
			slog.Error("Failed to query users", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var users []models.User
		users, err = pgx.CollectRows[models.User](rows, pgx.RowToStructByNameLax[models.User])
		if err != nil {
			slog.Error("Failed to collect rows", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(users)
		if err != nil {
			slog.Error("Failed to encode users to JSON", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	})
}

func DeleteUserByID(dbPool *pgxpool.Pool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ctx = r.Context()
		userID := r.PathValue("id")

		if userID == "" {
			slog.Error("User ID is empty")
			http.Error(w, "User ID is required", http.StatusBadRequest)
			return
		}

		err := queries.DeleteUserByID(ctx, dbPool, userID)
		if err != nil {
			slog.Error("Failed to delete user", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}
