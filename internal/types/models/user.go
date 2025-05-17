package models

import "time"

type User struct {
	ID        int       `json:"id"`
	ClerkID   string    `json:"clerk_id"`
	CreatedAt time.Time `json:"created_at"`
}
