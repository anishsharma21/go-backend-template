package setup

import (
	"fmt"
	"net/http"

	"github.com/anishsharma21/go-web-dev-template/internal"
	"github.com/anishsharma21/go-web-dev-template/internal/handlers"
	"github.com/anishsharma21/go-web-dev-template/internal/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

type routeConfig struct {
	Handler      http.Handler
	ApplyLogging bool
	ApplyJWT     bool
}

func Routes(dbPool *pgxpool.Pool) *http.ServeMux {
	mux := http.NewServeMux()

	routes := map[string]routeConfig{
		fmt.Sprintf("POST /%s/signup", internal.API_VERSION): {
			Handler:      handlers.AddNewUser(dbPool),
			ApplyLogging: true,
			ApplyJWT:     false,
		},
		fmt.Sprintf("GET /%s/users", internal.API_VERSION): {
			Handler:      handlers.GetUsers(dbPool),
			ApplyLogging: true,
			ApplyJWT:     true,
		},
		fmt.Sprintf("GET /%s/users/{clerk_user_id}", internal.API_VERSION): {
			Handler:      handlers.GetUserByClerkUserId(dbPool),
			ApplyLogging: true,
			ApplyJWT:     true,
		},
		fmt.Sprintf("DELETE /%s/users/{id}", internal.API_VERSION): {
			Handler:      handlers.DeleteUserByID(dbPool),
			ApplyLogging: true,
			ApplyJWT:     true,
		},

		"GET /static/": {
			Handler:      http.StripPrefix("/static/", http.FileServer(http.Dir("static"))),
			ApplyLogging: false,
			ApplyJWT:     false,
		},
	}

	for pattern, config := range routes {
		handler := config.Handler
		if config.ApplyLogging {
			handler = middleware.LoggingMiddleware(handler)
		}
		if config.ApplyJWT {
			handler = middleware.ClerkAuthMiddleware(handler)
		}
		mux.Handle(pattern, handler)
	}

	return mux
}
