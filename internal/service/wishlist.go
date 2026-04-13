package service

import (
	"context"
	"time"
	"wishlist-api/internal/models"
	"wishlist-api/internal/repository"

	"github.com/google/uuid"
)

// WishlistService handles wishlist operations.
type WishlistService struct {
	wishlistRepo repository.WishlistRepository
}

// NewWishlistService creates a new WishlistService.
func NewWishlistService(wishlistRepo repository.WishlistRepository) *WishlistService {
	return &WishlistService{wishlistRepo: wishlistRepo}
}

// Create creates a new wishlist for the user.
func (s *WishlistService) Create(ctx context.Context, userID int, title, description string, eventDate time.Time) (*models.Wishlist, error) {
	if eventDate.Before(time.Now()) {
		return nil, ErrInvalidInput
	}
	return s.wishlistRepo.Create(ctx, userID, title, description, eventDate)
}

// GetByID returns a wishlist by ID if it belongs to the user.
func (s *WishlistService) GetByID(ctx context.Context, id, userID int) (*models.Wishlist, error) {
	w, err := s.wishlistRepo.GetByID(ctx, id)
	if err != nil || w == nil {
		return nil, ErrNotFound
	}
	if w.UserID != userID {
		return nil, ErrForbidden
	}
	return w, nil
}

// GetByAccessToken returns a wishlist by its public access token.
func (s *WishlistService) GetByAccessToken(ctx context.Context, token uuid.UUID) (*models.Wishlist, error) {
	w, err := s.wishlistRepo.GetByAccessToken(ctx, token)
	if err != nil || w == nil {
		return nil, ErrNotFound
	}
	return w, nil
}

// GetAllByUser returns all wishlists belonging to the user.
func (s *WishlistService) GetAllByUser(ctx context.Context, userID int) ([]models.Wishlist, error) {
	return s.wishlistRepo.GetAllByUser(ctx, userID)
}

// Update modifies an existing wishlist.
func (s *WishlistService) Update(ctx context.Context, id, userID int, title, description string, eventDate time.Time) error {
	w, err := s.wishlistRepo.GetByID(ctx, id)
	if err != nil || w == nil {
		return ErrNotFound
	}
	if w.UserID != userID {
		return ErrForbidden
	}
	if eventDate.Before(time.Now()) {
		return ErrInvalidInput
	}
	w.Title = title
	w.Description = description
	w.EventDate = eventDate
	return s.wishlistRepo.Update(ctx, w)
}

// Delete removes a wishlist.
func (s *WishlistService) Delete(ctx context.Context, id, userID int) error {
	w, err := s.wishlistRepo.GetByID(ctx, id)
	if err != nil || w == nil {
		return ErrNotFound
	}
	if w.UserID != userID {
		return ErrForbidden
	}
	return s.wishlistRepo.Delete(ctx, id)
}
