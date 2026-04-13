// Package models contains domain entities for the application.
package models

import (
	"time"

	"github.com/google/uuid"
)

// Item represents a wishlist item (gift).
type Item struct {
	ID          uuid.UUID `json:"id"`
	WishlistID  uuid.UUID `json:"wishlist_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ProductLink string    `json:"product_link"`
	Priority    Priority  `json:"priority"`
	IsBooked    bool      `json:"is_booked"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
