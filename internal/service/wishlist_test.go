package service

import (
	"context"
	"errors"
	"testing"
	"time"
	"wishlist-api/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWishlistService_Create_Success(t *testing.T) {
	mockRepo := new(MockWishlistRepository)
	service := NewWishlistService(mockRepo)

	userID := 1
	title := "Birthday"
	desc := "Gifts"
	eventDate := time.Now().AddDate(0, 1, 0)

	expectedWishlist := &models.Wishlist{
		ID:          1,
		UserID:      userID,
		Title:       title,
		Description: desc,
		EventDate:   eventDate,
	}

	mockRepo.On("Create", mock.Anything, userID, title, desc, eventDate).
		Return(expectedWishlist, nil)

	wishlist, err := service.Create(context.Background(), userID, title, desc, eventDate)
	assert.NoError(t, err)
	assert.Equal(t, expectedWishlist, wishlist)
	mockRepo.AssertExpectations(t)
}

func TestWishlistService_Create_PastDate(t *testing.T) {
	mockRepo := new(MockWishlistRepository)
	service := NewWishlistService(mockRepo)

	eventDate := time.Now().AddDate(0, 0, -1)
	_, err := service.Create(context.Background(), 1, "Title", "", eventDate)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrInvalidInput))
}

func TestWishlistService_GetByID_Success(t *testing.T) {
	mockRepo := new(MockWishlistRepository)
	service := NewWishlistService(mockRepo)

	wishlist := &models.Wishlist{ID: 1, UserID: 1, Title: "Test"}
	mockRepo.On("GetByID", mock.Anything, 1).Return(wishlist, nil)

	result, err := service.GetByID(context.Background(), 1, 1)
	assert.NoError(t, err)
	assert.Equal(t, wishlist, result)
	mockRepo.AssertExpectations(t)
}

func TestWishlistService_GetByID_Forbidden(t *testing.T) {
	mockRepo := new(MockWishlistRepository)
	service := NewWishlistService(mockRepo)

	wishlist := &models.Wishlist{ID: 1, UserID: 2, Title: "Test"}
	mockRepo.On("GetByID", mock.Anything, 1).Return(wishlist, nil)

	_, err := service.GetByID(context.Background(), 1, 1)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrForbidden))
}

func TestWishlistService_Update_Success(t *testing.T) {
	mockRepo := new(MockWishlistRepository)
	service := NewWishlistService(mockRepo)

	wishlist := &models.Wishlist{ID: 1, UserID: 1, Title: "Old", Description: "Old desc", EventDate: time.Now().AddDate(0, 1, 0)}
	mockRepo.On("GetByID", mock.Anything, 1).Return(wishlist, nil)
	mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(w *models.Wishlist) bool {
		return w.Title == "New Title"
	})).Return(nil)

	err := service.Update(context.Background(), 1, 1, "New Title", "New desc", wishlist.EventDate)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestWishlistService_Delete_Success(t *testing.T) {
	mockRepo := new(MockWishlistRepository)
	service := NewWishlistService(mockRepo)

	wishlist := &models.Wishlist{ID: 1, UserID: 1}
	mockRepo.On("GetByID", mock.Anything, 1).Return(wishlist, nil)
	mockRepo.On("Delete", mock.Anything, 1).Return(nil)

	err := service.Delete(context.Background(), 1, 1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
