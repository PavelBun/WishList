package service

import (
	"context"
	"errors"
	"wishlist-api/internal/models"
	"wishlist-api/internal/repository"
)

// ItemService handles business logic for wishlist items.
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
		return nil, errors.New("priority must be between 1 and 5")
	}
	w, err := s.wishlistRepo.GetByID(ctx, wishlistID)
	if err != nil || w == nil {
		return nil, errors.New("wishlist not found")
	}
	if w.UserID != userID {
		return nil, errors.New("access denied")
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

// GetByID returns an item by its ID after verifying access.
func (s *ItemService) GetByID(ctx context.Context, id, userID int) (*models.Item, error) {
	item, err := s.itemRepo.GetByID(ctx, id)
	if err != nil || item == nil {
		return nil, errors.New("item not found")
	}
	w, err := s.wishlistRepo.GetByID(ctx, item.WishlistID)
	if err != nil || w == nil {
		return nil, errors.New("wishlist not found")
	}
	if w.UserID != userID {
		return nil, errors.New("access denied")
	}
	return item, nil
}

// GetAllByWishlistID returns all items of a wishlist after verifying access.
func (s *ItemService) GetAllByWishlistID(ctx context.Context, wishlistID, userID int) ([]models.Item, error) {
	w, err := s.wishlistRepo.GetByID(ctx, wishlistID)
	if err != nil || w == nil {
		return nil, errors.New("wishlist not found")
	}
	if w.UserID != userID {
		return nil, errors.New("access denied")
	}
	return s.itemRepo.GetAllByWishlistID(ctx, wishlistID)
}

// GetAllPublicByWishlistID returns all items of a wishlist without access check.
func (s *ItemService) GetAllPublicByWishlistID(ctx context.Context, wishlistID int) ([]models.Item, error) {
	return s.itemRepo.GetAllByWishlistID(ctx, wishlistID)
}

func (s *ItemService) getItemAndVerifyAccess(ctx context.Context, itemID, userID int) (*models.Item, *models.Wishlist, error) {
	item, err := s.itemRepo.GetByID(ctx, itemID)
	if err != nil || item == nil {
		return nil, nil, errors.New("item not found")
	}
	w, err := s.wishlistRepo.GetByID(ctx, item.WishlistID)
	if err != nil || w == nil {
		return nil, nil, errors.New("wishlist not found")
	}
	if w.UserID != userID {
		return nil, nil, errors.New("access denied")
	}
	return item, w, nil
}

// Update modifies an existing item.
func (s *ItemService) Update(ctx context.Context, id, userID int, title, description, productLink string, priority int) error {
	item, _, err := s.getItemAndVerifyAccess(ctx, id, userID)
	if err != nil {
		return err
	}
	item.Title = title
	item.Description = description
	item.ProductLink = productLink
	item.Priority = priority
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

// BookItem marks an item as reserved.
func (s *ItemService) BookItem(ctx context.Context, id int, wishlistID int) error {
	item, err := s.itemRepo.GetByID(ctx, id)
	if err != nil || item == nil {
		return errors.New("item not found")
	}
	if item.WishlistID != wishlistID {
		return errors.New("item does not belong to this wishlist")
	}
	return s.itemRepo.BookItem(ctx, id)
}
