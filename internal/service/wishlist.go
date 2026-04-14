package service

import (
	"context"
	"errors"
	"fmt"
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
func (s *WishlistService) Create(ctx context.Context, userID uuid.UUID, title, description string, eventDate time.Time) (*models.Wishlist, error) {
	today := time.Now().Truncate(24 * time.Hour).UTC()
	eventDateNormalized := eventDate.Truncate(24 * time.Hour).UTC()
	if eventDateNormalized.Before(today) {
		return nil, ErrInvalidInput
	}
	wishlist, err := s.wishlistRepo.Create(ctx, userID, title, description, eventDate)
	if err != nil {
		return nil, fmt.Errorf("failed to create wishlist: %w", err)
	}
	return wishlist, nil
}

// GetByID returns a wishlist by ID if it belongs to the user.
func (s *WishlistService) GetByID(ctx context.Context, id, userID uuid.UUID) (*models.Wishlist, error) {
	w, err := s.wishlistRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get wishlist: %w", err)
	}
	if w.UserID != userID {
		return nil, ErrForbidden
	}
	return w, nil
}

// GetByAccessToken returns a wishlist by its public access token.
func (s *WishlistService) GetByAccessToken(ctx context.Context, token uuid.UUID) (*models.Wishlist, error) {
	w, err := s.wishlistRepo.GetByAccessToken(ctx, token)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get wishlist by token: %w", err)
	}
	return w, nil
}

// GetAllByUser returns all wishlists belonging to the user.
func (s *WishlistService) GetAllByUser(ctx context.Context, userID uuid.UUID) ([]models.Wishlist, error) {
	wishlists, err := s.wishlistRepo.GetAllByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list wishlists: %w", err)
	}
	return wishlists, nil
}

// Update modifies an existing wishlist.
// Fields that are nil will not be updated.
func (s *WishlistService) Update(ctx context.Context, id, userID uuid.UUID, title *string, description *string, eventDate *time.Time) error {
	w, err := s.wishlistRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("failed to get wishlist for update: %w", err)
	}
	if w.UserID != userID {
		return ErrForbidden
	}

	if title != nil {
		w.Title = *title
	}
	if description != nil {
		w.Description = *description
	}
	if eventDate != nil {
		today := time.Now().Truncate(24 * time.Hour).UTC()
		eventDateNormalized := eventDate.Truncate(24 * time.Hour).UTC()
		if eventDateNormalized.Before(today) {
			return ErrInvalidInput
		}
		w.EventDate = *eventDate
	}

	if err := s.wishlistRepo.Update(ctx, w); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("failed to update wishlist: %w", err)
	}
	return nil
}

// Delete removes a wishlist.
func (s *WishlistService) Delete(ctx context.Context, id, userID uuid.UUID) error {
	w, err := s.wishlistRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("failed to get wishlist for delete: %w", err)
	}
	if w.UserID != userID {
		return ErrForbidden
	}
	if err := s.wishlistRepo.Delete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("failed to delete wishlist: %w", err)
	}
	return nil
}
