package handlers

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/anishsharma21/go-web-dev-template/internal"
	"github.com/anishsharma21/go-web-dev-template/internal/db"
	"github.com/anishsharma21/go-web-dev-template/internal/types/models"
	"github.com/jackc/pgx/v5/pgxpool"
	svix "github.com/svix/svix-webhooks/go"
)

// ClerkUserCreated represents the user.created event payload
type ClerkUserCreated struct {
	ID             string `json:"id"`
	Object         string `json:"object"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	EmailAddresses []struct {
		ID           string `json:"id"`
		EmailAddress string `json:"email_address"`
	} `json:"email_addresses"`
	// Add other fields as needed
}

var WEBHOOK_SECRET string

func init() {
	WEBHOOK_SECRET = os.Getenv(internal.CLERK_WEBHOOK_SIGNING_SECRET)
	if WEBHOOK_SECRET == "" {
		slog.ErrorContext(context.Background(), "CLERK_WEBHOOK_SIGNING_SECRET not set")
		os.Exit(1)
	}
}

// ClerkWebhookHandler handles webhook events from Clerk
func ClerkWebhookHandler(dbPool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Read the request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "Failed to read request body", slog.String("error", err.Error()))
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}

		// Initialize the Svix webhook verifier
		wh, err := svix.NewWebhook(WEBHOOK_SECRET)
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "Failed to initialize webhook verifier", slog.String("error", err.Error()))
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		// Get headers needed for verification
		headers := http.Header{}
		for _, header := range []string{"svix-id", "svix-timestamp", "svix-signature"} {
			if value := r.Header.Get(header); value != "" {
				headers.Set(header, value)
			} else {
				slog.LogAttrs(ctx, slog.LevelError, "Missing required header", slog.String("header", header))
				http.Error(w, "Missing required headers", http.StatusBadRequest)
				return
			}
		}

		// Verify the webhook
		err = wh.Verify(body, headers)
		if err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "Invalid webhook signature", slog.String("error", err.Error()))
			http.Error(w, "Invalid webhook signature", http.StatusUnauthorized)
			return
		}

		// Parse the webhook payload
		var payload struct {
			Type string          `json:"type"`
			Data json.RawMessage `json:"data"`
		}
		if err := json.Unmarshal(body, &payload); err != nil {
			slog.LogAttrs(ctx, slog.LevelError, "Failed to parse webhook payload", slog.String("error", err.Error()))
			http.Error(w, "Failed to parse webhook payload", http.StatusBadRequest)
			return
		}

		// Handle different event types
		slog.LogAttrs(ctx, slog.LevelInfo, "Received verified webhook", slog.String("event_type", payload.Type))

		switch payload.Type {
		case "user.created":
			var userData ClerkUserCreated
			if err := json.Unmarshal(payload.Data, &userData); err != nil {
				slog.LogAttrs(ctx, slog.LevelError, "Failed to parse user data", slog.String("error", err.Error()))
				http.Error(w, "Failed to parse user data", http.StatusBadRequest)
				return
			}

			// Log the received data to inspect the actual structure
			slog.LogAttrs(ctx, slog.LevelInfo, "Received user data",
				slog.String("user_id", userData.ID),
				slog.String("name", userData.FirstName+" "+userData.LastName))

			// Create a new user model from the webhook data
			user := models.User{
				ClerkID: userData.ID,
				// CreatedAt is handled by the database default value
			}

			// Add the user to the database
			if err := db.AddUser(ctx, dbPool, user); err != nil {
				slog.LogAttrs(ctx, slog.LevelError, "Failed to add user to database",
					slog.String("error", err.Error()),
					slog.String("clerk_id", userData.ID))
				http.Error(w, "Failed to process user data", http.StatusInternalServerError)
				return
			}

			slog.LogAttrs(ctx, slog.LevelInfo, "Successfully added user to database",
				slog.String("clerk_id", userData.ID))

		// Handle other event types as needed
		default:
			slog.LogAttrs(ctx, slog.LevelInfo, "Unhandled event type", slog.String("type", payload.Type))
		}

		// Return success
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Webhook received"))
	}
}
