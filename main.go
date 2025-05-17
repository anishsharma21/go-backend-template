package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/anishsharma21/go-web-dev-template/internal"
	"github.com/anishsharma21/go-web-dev-template/internal/middleware"
	"github.com/anishsharma21/go-web-dev-template/internal/setup"
	"github.com/jackc/pgx/v5/pgxpool"
)

var dbConnStr string
var dbPool *pgxpool.Pool

func init() {
	dbConnStr = os.Getenv(internal.DATABASE_URL)
	if dbConnStr == "" {
		slog.Error("DATABASE_URL environment variable not set")
		os.Exit(1)
	}

	defaultHandler := slog.NewJSONHandler(os.Stdout, nil)
	slog.SetDefault(slog.New(&middleware.CustomLogHandler{Handler: defaultHandler}))
}

func main() {
	// Setup context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbPool, err := setup.DBPool(ctx, dbConnStr)
	if err != nil {
		slog.Error("Failed to initialise database connection pool", "error", err)
		return
	}
	defer dbPool.Close()

	if os.Getenv(internal.RUN_MIGRATION) == "true" {
		slog.Info("Attempting to run database migrations...")
		err := setup.Migrations(dbConnStr)
		if err != nil {
			slog.Error("Failed to run database migrations", "error", err)
			return
		}
		slog.Info("Database migrations complete.")
	} else {
		slog.Info("Database migrations skipped.")
	}

	port := os.Getenv(internal.PORT)
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: setup.Routes(dbPool),
		BaseContext: func(l net.Listener) context.Context {
			url := "http://" + l.Addr().String()
			slog.Info(fmt.Sprintf("Server started on %s", url))
			return ctx
		},
	}

	shutdownChan := make(chan bool, 1)

	// Start server in a goroutine
	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			slog.Error("HTTP server closed early", "error", err)
		}
		slog.Info("Stopped server new connections.")
		shutdownChan <- true
	}()

	// Listen for OS signals (SIGINT, SIGTERM) to shutdown server gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan
	slog.Warn("Received signal", "signal", sig.String())

	// Shutdown server gracefully within 10 seconds
	shutdownCtx, shutdownRelease := context.WithTimeout(ctx, 10*time.Second)
	defer shutdownRelease()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("HTTP shutdown error occurred", "error", err)
	}
	<-shutdownChan
	close(shutdownChan)

	slog.Info("Graceful server shutdown complete.")
}
