package db

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/anishsharma21/go-web-dev-template/internal/types/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func AddUser(ctx context.Context, dbPool *pgxpool.Pool, user models.User) error {
	query := "INSERT INTO users (clerk_id) VALUES ($1)"

	ct, err := dbPool.Exec(ctx, query, user.ClerkID)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	if ct.RowsAffected() != 1 {
		return fmt.Errorf("unexpected number of rows affected: %d", ct.RowsAffected())
	}

	slog.InfoContext(ctx, "User signed up successfully",
		"clerk_id", user.ClerkID,
		"command_tag", ct.String())

	return nil
}

func GetUserByClerkUserId(ctx context.Context, dbPool *pgxpool.Pool, clerkUserId string) (models.User, error) {
	query := "SELECT * FROM users WHERE clerk_id = $1"

	row := dbPool.QueryRow(ctx, query, clerkUserId)

	var user models.User
	if err := row.Scan(&user.ID, &user.ClerkID, &user.CreatedAt); err != nil {
		if err == pgx.ErrNoRows {
			return models.User{}, err
		}
		return models.User{}, fmt.Errorf("error retrieving user with clerk_id %q: %w", clerkUserId, err)
	}
	return user, nil
}

func GetUsers(ctx context.Context, dbPool *pgxpool.Pool) ([]models.User, error) {
	query := "SELECT * FROM users"

	rows, err := dbPool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error retrieving users: %w", err)
	}
	defer rows.Close()

	var users []models.User
	users, err = pgx.CollectRows[models.User](rows, pgx.RowToStructByNameLax[models.User])
	if err != nil {
		return nil, fmt.Errorf("error collecting users: %w", err)
	}

	return users, nil
}

func DeleteUserByID(ctx context.Context, dbPool *pgxpool.Pool, id string) error {
	query := "DELETE FROM users WHERE id = $1"

	result, err := dbPool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("failed to delete user (no row affected): %w", err)
	}

	return nil
}
