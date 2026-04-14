package service

import (
	"context"
	"errors"
	"fmt"
	"wishlist-api/internal/models"
	"wishlist-api/internal/repository"

	"github.com/google/uuid"
)

// ItemService handles wishlist item operations.
type ItemService struct {
	itemRepo     repository.ItemRepository
	wishlistRepo repository.WishlistRepository
}

// NewItemService creates a new ItemService.
func NewItemService(itemRepo repository.ItemRepository, wishlistRepo repository.WishlistRepository) *ItemService {
	return &ItemService{
		itemRepo:     itemRepo,
		wishlistRepo: wishlistRepo,
	}
}

// Create adds a new item to a wishlist.
func (s *ItemService) Create(ctx context.Context, wishlistID, userID uuid.UUID, title, description, productLink string, priority models.Priority) (*models.Item, error) {
	if !priority.Valid() {
		return nil, ErrInvalidInput
	}
	w, err := s.wishlistRepo.GetByID(ctx, wishlistID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get wishlist for item creation: %w", err)
	}
	if w.UserID != userID {
		return nil, ErrForbidden
	}
	item := &models.Item{
		WishlistID:  wishlistID,
		Title:       title,
		Description: description,
		ProductLink: productLink,
		Priority:    priority,
	}
	if err := s.itemRepo.Create(ctx, item); err != nil {
		return nil, fmt.Errorf("failed to create item: %w", err)
	}
	return item, nil
}

// GetByID returns an item by ID, checking user ownership.
func (s *ItemService) GetByID(ctx context.Context, id, userID uuid.UUID) (*models.Item, error) {
	item, err := s.itemRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get item: %w", err)
	}
	w, err := s.wishlistRepo.GetByID(ctx, item.WishlistID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get wishlist for item: %w", err)
	}
	if w.UserID != userID {
		return nil, ErrForbidden
	}
	return item, nil
}

// GetAllByWishlistID returns all items of a wishlist, checking user ownership.
func (s *ItemService) GetAllByWishlistID(ctx context.Context, wishlistID, userID uuid.UUID) ([]models.Item, error) {
	w, err := s.wishlistRepo.GetByID(ctx, wishlistID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get wishlist for items listing: %w", err)
	}
	if w.UserID != userID {
		return nil, ErrForbidden
	}
	items, err := s.itemRepo.GetAllByWishlistID(ctx, wishlistID)
	if err != nil {
		return nil, fmt.Errorf("failed to list items: %w", err)
	}
	if items == nil {
		return []models.Item{}, nil
	}
	return items, nil
}

// GetAllPublicByWishlistID returns all items of a wishlist without auth checks.
func (s *ItemService) GetAllPublicByWishlistID(ctx context.Context, wishlistID uuid.UUID) ([]models.Item, error) {
	items, err := s.itemRepo.GetAllByWishlistID(ctx, wishlistID)
	if err != nil {
		return nil, fmt.Errorf("failed to list public items: %w", err)
	}
	if items == nil {
		return []models.Item{}, nil
	}
	return items, nil
}

// Update modifies an existing item.
func (s *ItemService) Update(ctx context.Context, id, userID uuid.UUID, title *string, description *string, productLink *string, priority *models.Priority) error {
	if priority != nil && !priority.Valid() {
		return ErrInvalidInput
	}
	item, _, err := s.getItemAndVerifyAccess(ctx, id, userID)
	if err != nil {
		return err
	}
	if title != nil {
		item.Title = *title
	}
	if description != nil {
		item.Description = *description
	}
	if productLink != nil {
		item.ProductLink = *productLink
	}
	if priority != nil {
		item.Priority = *priority
	}
	if err := s.itemRepo.Update(ctx, item); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("failed to update item: %w", err)
	}
	return nil
}

// Delete removes an item.
func (s *ItemService) Delete(ctx context.Context, id, userID uuid.UUID) error {
	_, _, err := s.getItemAndVerifyAccess(ctx, id, userID)
	if err != nil {
		return err
	}
	if err := s.itemRepo.Delete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("failed to delete item: %w", err)
	}
	return nil
}

// BookItem marks an item as booked. Returns ErrAlreadyBooked if already booked.
func (s *ItemService) BookItem(ctx context.Context, id, wishlistID uuid.UUID) error {
	item, err := s.itemRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("failed to get item for booking: %w", err)
	}
	if item.WishlistID != wishlistID {
		return ErrForbidden
	}
	err = s.itemRepo.BookItem(ctx, id)
	if errors.Is(err, repository.ErrAlreadyBooked) {
		return ErrAlreadyBooked
	}
	if errors.Is(err, repository.ErrNotFound) {
		return ErrNotFound
	}
	if err != nil {
		return fmt.Errorf("failed to book item: %w", err)
	}
	return nil
}

// getItemAndVerifyAccess retrieves an item and its wishlist, checking user access.
func (s *ItemService) getItemAndVerifyAccess(ctx context.Context, itemID, userID uuid.UUID) (*models.Item, *models.Wishlist, error) {
	item, err := s.itemRepo.GetByID(ctx, itemID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, nil, ErrNotFound
		}
		return nil, nil, fmt.Errorf("failed to get item for access check: %w", err)
	}
	w, err := s.wishlistRepo.GetByID(ctx, item.WishlistID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, nil, ErrNotFound
		}
		return nil, nil, fmt.Errorf("failed to get wishlist for access check: %w", err)
	}
	if w.UserID != userID {
		return nil, nil, ErrForbidden
	}
	return item, w, nil
}
