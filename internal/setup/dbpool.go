package setup

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func DBPool(ctx context.Context, dbConnStr string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(dbConnStr)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse database connection string: %v", err)
	}

	// Set connection pool configurations
	// Sets the maximum time an idle connection can remain in the pool before being closed
	config.MaxConnIdleTime = 1 * time.Minute
	// To prevent database and backend from ever sleeping, uncomment the following line
	config.MinConns = 1

	var dbPool *pgxpool.Pool
	for i := 1; i <= 5; i++ {
		dbPool, err = pgxpool.NewWithConfig(ctx, config)
		if err == nil && dbPool != nil {
			break
		}
		slog.Warn("Failed to initialise database connection pool", "error", err)
		slog.Info(fmt.Sprintf("Retrying in %d seconds...", i*i))
		time.Sleep(time.Duration(i*i) * time.Second)
	}
	if dbPool == nil {
		return nil, fmt.Errorf("Failed to initialise database connection pool after 5 attempts")
	}

	for i := 1; i <= 5; i++ {
		err = dbPool.Ping(ctx)
		if err == nil && dbPool != nil {
			break
		}
		slog.Warn("Failed to ping database connection pool", "error", err)
		slog.Info(fmt.Sprintf("Retrying in %d seconds...", i*i))
		time.Sleep(time.Duration(i*i) * time.Second)
	}
	if err != nil {
		return nil, fmt.Errorf("Failed to ping database connection pool after 5 attempts")
	}

	return dbPool, nil
}
