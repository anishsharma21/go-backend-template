package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/anishsharma21/go-web-dev-template/internal"
	"github.com/anishsharma21/go-web-dev-template/internal/handlers"
	"github.com/jackc/pgx/v5/pgxpool"
)

var dbPool *pgxpool.Pool
var ctx, cancel = context.WithCancel(context.Background())

func TestMain(m *testing.M) {
	dbConnStr := os.Getenv(internal.DATABASE_URL)
	config, err := pgxpool.ParseConfig(dbConnStr)
	if err != nil {
		log.Fatalf("Failed to parse database connection string.\n")
	}

	for i := 1; i <= 5; i++ {
		dbPool, err = pgxpool.NewWithConfig(ctx, config)
		if err == nil && dbPool != nil {
			break
		}
		log.Printf("Failed to initialise database connection pool")
		log.Printf(fmt.Sprintf("Retrying in %d seconds...", i*i))
		time.Sleep(time.Duration(i*i) * time.Second)
	}
	if dbPool == nil {
		log.Fatalf("Failed to initialise database connection pool after 5 attempts")
	}
	defer dbPool.Close()

	// Run the tests
	code := m.Run()

	cancel()

	os.Exit(code)
}

func TestUserSignUpFlow(t *testing.T) {
	// Arrange
	ts := httptest.NewServer(handlers.AddNewUser(dbPool))
	defer ts.Close()

	clerkID := "testclerkid"

	userData := map[string]string{
		"clerk_id": clerkID,
	}
	jsonData, err := json.Marshal(userData)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	// Act
	resp, err := ts.Client().Post(ts.URL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Expected no error when sending POST request, got %v\n", err)
	}
	defer resp.Body.Close()

	// Assert
	if resp.StatusCode != 201 {
		t.Errorf("Expected status code 201, got %v\n", resp.StatusCode)
	}

	// Teardown
	_, err = dbPool.Exec(ctx, "DELETE FROM users WHERE clerk_id = $1", clerkID)
	if err != nil {
		t.Fatalf("Failed to delete user from database, %v\n", err)
	}
}
