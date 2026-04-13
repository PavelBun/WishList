// Package service contains business logic for the application.
package service

import (
	"context"
	"errors"
	"time"
	"wishlist-api/internal/models"
	"wishlist-api/internal/repository"

	"github.com/google/uuid"
)

// WishlistService handles business logic for wishlists.
type WishlistService struct {
	wishlistRepo repository.WishlistRepository
}

// NewWishlistService creates a new WishlistService.
func NewWishlistService(wishlistRepo repository.WishlistRepository) *WishlistService {
	return &WishlistService{wishlistRepo: wishlistRepo}
}

// Create creates a new wishlist for a user.
func (s *WishlistService) Create(ctx context.Context, userID int, title, description string, eventDate time.Time) (*models.Wishlist, error) {
	if eventDate.Before(time.Now()) {
		return nil, errors.New("event date must be in the future")
	}
	return s.wishlistRepo.Create(ctx, userID, title, description, eventDate)
}

// GetByID returns a wishlist by its ID, ensuring the user has access.
func (s *WishlistService) GetByID(ctx context.Context, id, userID int) (*models.Wishlist, error) {
	w, err := s.wishlistRepo.GetByID(ctx, id)
	if err != nil || w == nil {
		return nil, errors.New("wishlist not found")
	}
	if w.UserID != userID {
		return nil, errors.New("access denied")
	}
	return w, nil
}

// GetByAccessToken returns a wishlist by its public access token.
func (s *WishlistService) GetByAccessToken(ctx context.Context, token uuid.UUID) (*models.Wishlist, error) {
	return s.wishlistRepo.GetByAccessToken(ctx, token)
}

// GetAllByUser returns all wishlists for the given user.
func (s *WishlistService) GetAllByUser(ctx context.Context, userID int) ([]models.Wishlist, error) {
	return s.wishlistRepo.GetAllByUser(ctx, userID)
}

// Update updates an existing wishlist if the user owns it.
func (s *WishlistService) Update(ctx context.Context, id, userID int, title, description string, eventDate time.Time) error {
	w, err := s.wishlistRepo.GetByID(ctx, id)
	if err != nil || w == nil {
		return errors.New("wishlist not found")
	}
	if w.UserID != userID {
		return errors.New("access denied")
	}
	if eventDate.Before(time.Now().Truncate(24 * time.Hour)) {
		return errors.New("event date must be in the future")
	}
	w.Title = title
	w.Description = description
	w.EventDate = eventDate
	return s.wishlistRepo.Update(ctx, w)
}

// Delete removes a wishlist if the user owns it.
func (s *WishlistService) Delete(ctx context.Context, id, userID int) error {
	w, err := s.wishlistRepo.GetByID(ctx, id)
	if err != nil || w == nil {
		return errors.New("wishlist not found")
	}
	if w.UserID != userID {
		return errors.New("access denied")
	}
	return s.wishlistRepo.Delete(ctx, id)
}
