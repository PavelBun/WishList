package dto

// CreateWishlistRequest represents the request body for creating a wishlist.
type CreateWishlistRequest struct {
	Title       string `json:"title" validate:"required,min=1,max=255"`
	Description string `json:"description" validate:"max=1000"`
	EventDate   string `json:"event_date" validate:"required"`
}

// UpdateWishlistRequest represents the request body for updating a wishlist.
type UpdateWishlistRequest struct {
	Title       *string `json:"title" validate:"omitempty,min=1,max=255"`
	Description *string `json:"description" validate:"omitempty,max=1000"`
	EventDate   *string `json:"event_date" validate:"omitempty"`
}
