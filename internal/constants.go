package internal

import (
	"log/slog"
	"os"
)

// Context keys
const CLERK_USER_ID_KEY = "clerk_user_id"
const REQUEST_ID_KEY = "request_id"

// Environment variable keys
const CLERK_SECRET_KEY = "CLERK_SECRET_KEY"
const CLERK_WEBHOOK_SIGNING_SECRET = "CLERK_WEBHOOK_SIGNING_SECRET"
const DATABASE_URL = "DATABASE_URL"
const GOOSE_DRIVER = "GOOSE_DRIVER"
const GOOSE_DBSTRING = "GOOSE_DBSTRING"
const GOOSE_MIGRATION_DIR = "GOOSE_MIGRATION_DIR"
const RUN_MIGRATION = "RUN_MIGRATION"
const PORT = "PORT"

var ENVIRONMENT string

// API version
const API_VERSION = "v1"

func init() {
	ENVIRONMENT = os.Getenv("ENVIRONMENT")
	if ENVIRONMENT == "" {
		slog.Warn("ENVIRONMENT not set, defaulting to development")
		ENVIRONMENT = "development"
	}
}
