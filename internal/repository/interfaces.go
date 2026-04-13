// Package repository defines interfaces for data access.
package repository

import (
	"context"
	"time"
	"wishlist-api/internal/models"

	"github.com/google/uuid"
)

// UserRepository defines methods for user persistence.
type UserRepository interface {
	Create(ctx context.Context, email, passwordHash string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
}

// WishlistRepository defines methods for wishlist persistence.
type WishlistRepository interface {
	Create(ctx context.Context, userID int, title, description string, eventDate time.Time) (*models.Wishlist, error)
	GetByID(ctx context.Context, id int) (*models.Wishlist, error)
	GetByAccessToken(ctx context.Context, token uuid.UUID) (*models.Wishlist, error)
	GetAllByUser(ctx context.Context, userID int) ([]models.Wishlist, error)
	Update(ctx context.Context, w *models.Wishlist) error
	Delete(ctx context.Context, id int) error
}

// ItemRepository defines methods for item persistence.
type ItemRepository interface {
	Create(ctx context.Context, item *models.Item) error
	GetByID(ctx context.Context, id int) (*models.Item, error)
	GetAllByWishlistID(ctx context.Context, wishlistID int) ([]models.Item, error)
	Update(ctx context.Context, item *models.Item) error
	Delete(ctx context.Context, id int) error
	BookItem(ctx context.Context, id int) error
}
