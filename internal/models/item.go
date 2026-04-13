// Package models contains domain models for the application.
package models

import "time"

// Item represents a wishlist item (gift).
type Item struct {
	ID          int       `json:"id"`
	WishlistID  int       `json:"wishlist_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ProductLink string    `json:"product_link"`
	Priority    int       `json:"priority"`
	IsBooked    bool      `json:"is_booked"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
