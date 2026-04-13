package service

import (
	"context"
	"errors"
	"wishlist-api/internal/models"
	"wishlist-api/internal/repository"
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
func (s *ItemService) Create(ctx context.Context, wishlistID, userID int, title, description, productLink string, priority int) (*models.Item, error) {
	if priority < 1 || priority > 5 {
		return nil, ErrInvalidInput
	}
	w, err := s.wishlistRepo.GetByID(ctx, wishlistID)
	if err != nil || w == nil {
		return nil, ErrNotFound
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
	err = s.itemRepo.Create(ctx, item)
	if err != nil {
		return nil, err
	}
	return item, nil
}

// GetByID returns an item by ID, checking user ownership.
func (s *ItemService) GetByID(ctx context.Context, id, userID int) (*models.Item, error) {
	item, err := s.itemRepo.GetByID(ctx, id)
	if err != nil || item == nil {
		return nil, ErrNotFound
	}
	w, err := s.wishlistRepo.GetByID(ctx, item.WishlistID)
	if err != nil || w == nil {
		return nil, ErrNotFound
	}
	if w.UserID != userID {
		return nil, ErrForbidden
	}
	return item, nil
}

// GetAllByWishlistID returns all items of a wishlist, checking user ownership.
func (s *ItemService) GetAllByWishlistID(ctx context.Context, wishlistID, userID int) ([]models.Item, error) {
	w, err := s.wishlistRepo.GetByID(ctx, wishlistID)
	if err != nil || w == nil {
		return nil, ErrNotFound
	}
	if w.UserID != userID {
		return nil, ErrForbidden
	}
	return s.itemRepo.GetAllByWishlistID(ctx, wishlistID)
}

// GetAllPublicByWishlistID returns all items of a wishlist without auth checks.
func (s *ItemService) GetAllPublicByWishlistID(ctx context.Context, wishlistID int) ([]models.Item, error) {
	return s.itemRepo.GetAllByWishlistID(ctx, wishlistID)
}

// Update modifies an existing item.
func (s *ItemService) Update(ctx context.Context, id, userID int, title *string, description *string, productLink *string, priority *int) error {
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
	return s.itemRepo.Update(ctx, item)
}

// Delete removes an item.
func (s *ItemService) Delete(ctx context.Context, id, userID int) error {
	_, _, err := s.getItemAndVerifyAccess(ctx, id, userID)
	if err != nil {
		return err
	}
	return s.itemRepo.Delete(ctx, id)
}

// BookItem marks an item as booked. Returns ErrAlreadyBooked if already booked.
func (s *ItemService) BookItem(ctx context.Context, id int, wishlistID int) error {
	item, err := s.itemRepo.GetByID(ctx, id)
	if err != nil || item == nil {
		return ErrNotFound
	}
	if item.WishlistID != wishlistID {
		return ErrForbidden
	}
	err = s.itemRepo.BookItem(ctx, id)
	if errors.Is(err, repository.ErrAlreadyBooked) {
		return ErrAlreadyBooked
	}
	return err
}

// getItemAndVerifyAccess retrieves an item and its wishlist, checking user access.
func (s *ItemService) getItemAndVerifyAccess(ctx context.Context, itemID, userID int) (*models.Item, *models.Wishlist, error) {
	item, err := s.itemRepo.GetByID(ctx, itemID)
	if err != nil || item == nil {
		return nil, nil, ErrNotFound
	}
	w, err := s.wishlistRepo.GetByID(ctx, item.WishlistID)
	if err != nil || w == nil {
		return nil, nil, ErrNotFound
	}
	if w.UserID != userID {
		return nil, nil, ErrForbidden
	}
	return item, w, nil
}
