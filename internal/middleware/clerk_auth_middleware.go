package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/anishsharma21/go-web-dev-template/internal"
	"github.com/clerk/clerk-sdk-go/v2"
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
)

func init() {
	apiKey := os.Getenv(internal.CLERK_SECRET_KEY)
	if apiKey == "" {
		slog.Error("CLERK_SECRET_KEY not set")
		os.Exit(1)
	}
	clerk.SetKey(apiKey)
}

// ClerkAuthMiddleware verifies JWT tokens and adds user ID to context
func ClerkAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		validateHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := clerk.SessionClaimsFromContext(r.Context())
			if !ok {
				slog.LogAttrs(ctx, slog.LevelError, "Failed to get session claims")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Extract user ID from claims
			userID := claims.Subject
			if userID == "" {
				slog.LogAttrs(ctx, slog.LevelError, "User ID not found in claims")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), internal.CLERK_USER_ID_KEY, userID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})

		clerkMiddleware := clerkhttp.RequireHeaderAuthorization()
		clerkHandler := clerkMiddleware(validateHandler)
		clerkHandler.ServeHTTP(w, r)
	})
}
