package models

import (
	"time"

	"github.com/google/uuid"
)

// Wishlist represents a user's wishlist for an event.
type Wishlist struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	EventDate   time.Time `json:"event_date"`
	AccessToken uuid.UUID `json:"access_token"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Items       []Item    `json:"items,omitempty"`
}
