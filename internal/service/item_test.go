package service

import (
	"context"
	"errors"
	"testing"
	"wishlist-api/internal/models"
	"wishlist-api/internal/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestItemService_Create_Success(t *testing.T) {
	mockItemRepo := new(MockItemRepository)
	mockWishlistRepo := new(MockWishlistRepository)
	service := NewItemService(mockItemRepo, mockWishlistRepo)

	wishlistID := 10
	userID := 1
	title := "PS5"
	desc := "Game console"
	link := "https://example.com"
	priority := 5

	wishlist := &models.Wishlist{ID: wishlistID, UserID: userID}
	mockWishlistRepo.On("GetByID", mock.Anything, wishlistID).Return(wishlist, nil)
	mockItemRepo.On("Create", mock.Anything, mock.MatchedBy(func(item *models.Item) bool {
		return item.WishlistID == wishlistID && item.Title == title
	})).Return(nil)

	item, err := service.Create(context.Background(), wishlistID, userID, title, desc, link, priority)
	assert.NoError(t, err)
	assert.NotNil(t, item)
	assert.Equal(t, title, item.Title)
	mockWishlistRepo.AssertExpectations(t)
	mockItemRepo.AssertExpectations(t)
}

func TestItemService_Create_InvalidPriority(t *testing.T) {
	mockItemRepo := new(MockItemRepository)
	mockWishlistRepo := new(MockWishlistRepository)
	service := NewItemService(mockItemRepo, mockWishlistRepo)

	_, err := service.Create(context.Background(), 1, 1, "title", "", "", 0)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrInvalidInput))
}

func TestItemService_BookItem_Success(t *testing.T) {
	mockItemRepo := new(MockItemRepository)
	mockWishlistRepo := new(MockWishlistRepository)
	service := NewItemService(mockItemRepo, mockWishlistRepo)

	itemID := 5
	wishlistID := 10
	item := &models.Item{ID: itemID, WishlistID: wishlistID}

	mockItemRepo.On("GetByID", mock.Anything, itemID).Return(item, nil)
	mockItemRepo.On("BookItem", mock.Anything, itemID).Return(nil)

	err := service.BookItem(context.Background(), itemID, wishlistID)
	assert.NoError(t, err)
	mockItemRepo.AssertExpectations(t)
}

func TestItemService_BookItem_AlreadyBooked(t *testing.T) {
	mockItemRepo := new(MockItemRepository)
	mockWishlistRepo := new(MockWishlistRepository)
	service := NewItemService(mockItemRepo, mockWishlistRepo)

	itemID := 5
	wishlistID := 10
	item := &models.Item{ID: itemID, WishlistID: wishlistID}

	mockItemRepo.On("GetByID", mock.Anything, itemID).Return(item, nil)
	mockItemRepo.On("BookItem", mock.Anything, itemID).Return(repository.ErrAlreadyBooked)

	err := service.BookItem(context.Background(), itemID, wishlistID)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrAlreadyBooked))
}

func TestItemService_GetByID_Success(t *testing.T) {
	mockItemRepo := new(MockItemRepository)
	mockWishlistRepo := new(MockWishlistRepository)
	service := NewItemService(mockItemRepo, mockWishlistRepo)

	item := &models.Item{ID: 1, WishlistID: 10}
	wishlist := &models.Wishlist{ID: 10, UserID: 1}

	mockItemRepo.On("GetByID", mock.Anything, 1).Return(item, nil)
	mockWishlistRepo.On("GetByID", mock.Anything, 10).Return(wishlist, nil)

	result, err := service.GetByID(context.Background(), 1, 1)
	assert.NoError(t, err)
	assert.Equal(t, item, result)
}

func TestItemService_Update_Success(t *testing.T) {
	mockItemRepo := new(MockItemRepository)
	mockWishlistRepo := new(MockWishlistRepository)
	service := NewItemService(mockItemRepo, mockWishlistRepo)

	item := &models.Item{ID: 1, WishlistID: 10, Title: "Old"}
	wishlist := &models.Wishlist{ID: 10, UserID: 1}

	mockItemRepo.On("GetByID", mock.Anything, 1).Return(item, nil)
	mockWishlistRepo.On("GetByID", mock.Anything, 10).Return(wishlist, nil)

	newTitle := "New Title"
	newPriority := 3
	mockItemRepo.On("Update", mock.Anything, mock.MatchedBy(func(i *models.Item) bool {
		return i.Title == newTitle && i.Priority == newPriority
	})).Return(nil)

	err := service.Update(context.Background(), 1, 1, &newTitle, nil, nil, &newPriority)
	assert.NoError(t, err)
	mockItemRepo.AssertExpectations(t)
}
