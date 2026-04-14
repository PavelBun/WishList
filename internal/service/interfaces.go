package service

import (
	"context"
	"time"
	"wishlist-api/internal/models"

	"github.com/google/uuid"
)

// AuthServiceInterface defines the methods for authentication operations.
type AuthServiceInterface interface {
	Register(ctx context.Context, email, password string) (*models.User, error)
	Login(ctx context.Context, email, password string) (string, error)
	ValidateToken(tokenString string) (uuid.UUID, error)
}

// WishlistServiceInterface defines the methods for wishlist operations.
type WishlistServiceInterface interface {
	Create(ctx context.Context, userID uuid.UUID, title, description string, eventDate time.Time) (*models.Wishlist, error)
	GetByID(ctx context.Context, id, userID uuid.UUID) (*models.Wishlist, error)
	GetByAccessToken(ctx context.Context, token uuid.UUID) (*models.Wishlist, error)
	GetAllByUser(ctx context.Context, userID uuid.UUID) ([]models.Wishlist, error)
	Update(ctx context.Context, id, userID uuid.UUID, title *string, description *string, eventDate *time.Time) error
	Delete(ctx context.Context, id, userID uuid.UUID) error
}

// ItemServiceInterface defines the methods for item operations.
type ItemServiceInterface interface {
	Create(ctx context.Context, wishlistID, userID uuid.UUID, title, description, productLink string, priority models.Priority) (*models.Item, error)
	GetByID(ctx context.Context, id, userID uuid.UUID) (*models.Item, error)
	GetAllByWishlistID(ctx context.Context, wishlistID, userID uuid.UUID) ([]models.Item, error)
	GetAllPublicByWishlistID(ctx context.Context, wishlistID uuid.UUID) ([]models.Item, error)
	Update(ctx context.Context, id, userID uuid.UUID, title *string, description *string, productLink *string, priority *models.Priority) error
	Delete(ctx context.Context, id, userID uuid.UUID) error
	BookItem(ctx context.Context, id, wishlistID uuid.UUID) error
}
