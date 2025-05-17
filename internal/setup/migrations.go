package setup

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/anishsharma21/go-web-dev-template/internal"
	"github.com/pressly/goose/v3"
)

func Migrations(dbConnStr string) error {
	if gooseDriver := os.Getenv(internal.GOOSE_DRIVER); gooseDriver == "" {
		return fmt.Errorf("Goose driver not set: GOOSE_DRIVER=?")
	}

	if gooseDbString := os.Getenv(internal.GOOSE_DBSTRING); gooseDbString == "" {
		return fmt.Errorf("Goose db string not set: GOOSE_DBSTRING=?")
	}

	if gooseMigrationDir := os.Getenv(internal.GOOSE_MIGRATION_DIR); gooseMigrationDir == "" {
		return fmt.Errorf("Goose migration dir not set: GOOSE_MIGRATION_DIR=?")
	}

	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		return fmt.Errorf("Failed to open database connection for *sql.DB: %v\n", err)
	}
	defer db.Close()

	if err = goose.Status(db, "migrations"); err != nil {
		return fmt.Errorf("Failed to retrieve status of migrations: %v\n", err)
	}

	if err = goose.Up(db, "migrations"); err != nil {
		return fmt.Errorf("Failed to run `goose up` command: %v\n", err)
	}

	return nil
}
