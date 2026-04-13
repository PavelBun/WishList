package service

import (
	"context"
	"errors"
	"testing"
	"wishlist-api/internal/models"
	"wishlist-api/internal/repository"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestItemService_Create_Success(t *testing.T) {
	mockItemRepo := new(MockItemRepository)
	mockWishlistRepo := new(MockWishlistRepository)
	svc := NewItemService(mockItemRepo, mockWishlistRepo)

	wishlistID := uuid.New()
	userID := uuid.New()
	title := "PS5"
	desc := "Game console"
	link := "https://example.com"
	priority := models.PriorityMust

	wishlist := &models.Wishlist{ID: wishlistID, UserID: userID}
	mockWishlistRepo.On("GetByID", mock.Anything, wishlistID).Return(wishlist, nil)
	mockItemRepo.On("Create", mock.Anything, mock.MatchedBy(func(item *models.Item) bool {
		return item.WishlistID == wishlistID && item.Title == title
	})).Return(nil)

	item, err := svc.Create(context.Background(), wishlistID, userID, title, desc, link, priority)
	assert.NoError(t, err)
	assert.NotNil(t, item)
	assert.Equal(t, title, item.Title)
	mockWishlistRepo.AssertExpectations(t)
	mockItemRepo.AssertExpectations(t)
}

func TestItemService_Create_InvalidPriority(t *testing.T) {
	mockItemRepo := new(MockItemRepository)
	mockWishlistRepo := new(MockWishlistRepository)
	svc := NewItemService(mockItemRepo, mockWishlistRepo)

	_, err := svc.Create(context.Background(), uuid.New(), uuid.New(), "title", "", "", models.Priority(0))
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrInvalidInput))
}

func TestItemService_BookItem_Success(t *testing.T) {
	mockItemRepo := new(MockItemRepository)
	mockWishlistRepo := new(MockWishlistRepository)
	svc := NewItemService(mockItemRepo, mockWishlistRepo)

	itemID := uuid.New()
	wishlistID := uuid.New()
	item := &models.Item{ID: itemID, WishlistID: wishlistID}

	mockItemRepo.On("GetByID", mock.Anything, itemID).Return(item, nil)
	mockItemRepo.On("BookItem", mock.Anything, itemID).Return(nil)

	err := svc.BookItem(context.Background(), itemID, wishlistID)
	assert.NoError(t, err)
	mockItemRepo.AssertExpectations(t)
}

func TestItemService_BookItem_AlreadyBooked(t *testing.T) {
	mockItemRepo := new(MockItemRepository)
	mockWishlistRepo := new(MockWishlistRepository)
	svc := NewItemService(mockItemRepo, mockWishlistRepo)

	itemID := uuid.New()
	wishlistID := uuid.New()
	item := &models.Item{ID: itemID, WishlistID: wishlistID}

	mockItemRepo.On("GetByID", mock.Anything, itemID).Return(item, nil)
	mockItemRepo.On("BookItem", mock.Anything, itemID).Return(repository.ErrAlreadyBooked)

	err := svc.BookItem(context.Background(), itemID, wishlistID)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrAlreadyBooked))
}

func TestItemService_GetByID_Success(t *testing.T) {
	mockItemRepo := new(MockItemRepository)
	mockWishlistRepo := new(MockWishlistRepository)
	svc := NewItemService(mockItemRepo, mockWishlistRepo)

	itemID := uuid.New()
	wishlistID := uuid.New()
	userID := uuid.New()
	item := &models.Item{ID: itemID, WishlistID: wishlistID}
	wishlist := &models.Wishlist{ID: wishlistID, UserID: userID}

	mockItemRepo.On("GetByID", mock.Anything, itemID).Return(item, nil)
	mockWishlistRepo.On("GetByID", mock.Anything, wishlistID).Return(wishlist, nil)

	result, err := svc.GetByID(context.Background(), itemID, userID)
	assert.NoError(t, err)
	assert.Equal(t, item, result)
}

func TestItemService_Update_Success(t *testing.T) {
	mockItemRepo := new(MockItemRepository)
	mockWishlistRepo := new(MockWishlistRepository)
	svc := NewItemService(mockItemRepo, mockWishlistRepo)

	itemID := uuid.New()
	wishlistID := uuid.New()
	userID := uuid.New()
	item := &models.Item{ID: itemID, WishlistID: wishlistID, Title: "Old", Priority: models.PriorityMedium}
	wishlist := &models.Wishlist{ID: wishlistID, UserID: userID}

	mockItemRepo.On("GetByID", mock.Anything, itemID).Return(item, nil)
	mockWishlistRepo.On("GetByID", mock.Anything, wishlistID).Return(wishlist, nil)

	newTitle := "New Title"
	newPriority := models.PriorityMust
	mockItemRepo.On("Update", mock.Anything, mock.MatchedBy(func(i *models.Item) bool {
		return i.Title == newTitle && i.Priority == newPriority
	})).Return(nil)

	err := svc.Update(context.Background(), itemID, userID, &newTitle, nil, nil, &newPriority)
	assert.NoError(t, err)
	mockItemRepo.AssertExpectations(t)
}
