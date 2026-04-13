package dto

// CreateItemRequest represents the request body for creating a wishlist item.
type CreateItemRequest struct {
	Title       string `json:"title" validate:"required,min=1,max=255"`
	Description string `json:"description" validate:"max=1000"`
	ProductLink string `json:"product_link" validate:"omitempty,url"`
	Priority    int    `json:"priority" validate:"required,min=1,max=5"`
}

// UpdateItemRequest represents the request body for updating a wishlist item.
type UpdateItemRequest struct {
	Title       string `json:"title" validate:"omitempty,min=1,max=255"`
	Description string `json:"description" validate:"max=1000"`
	ProductLink string `json:"product_link" validate:"omitempty,url"`
	Priority    int    `json:"priority" validate:"omitempty,min=1,max=5"`
}
