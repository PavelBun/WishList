package service

import (
	"context"
	"errors"
	"testing"
	"time"
	"wishlist-api/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWishlistService_Create_Success(t *testing.T) {
	mockRepo := new(MockWishlistRepository)
	service := NewWishlistService(mockRepo)

	userID := uuid.New()
	title := "Birthday"
	desc := "Gifts"
	eventDate := time.Now().AddDate(0, 1, 0)

	expectedWishlist := &models.Wishlist{
		ID:          uuid.New(),
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
	_, err := service.Create(context.Background(), uuid.New(), "Title", "", eventDate)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrInvalidInput))
}

func TestWishlistService_GetByID_Success(t *testing.T) {
	mockRepo := new(MockWishlistRepository)
	service := NewWishlistService(mockRepo)

	wishlistID := uuid.New()
	userID := uuid.New()
	wishlist := &models.Wishlist{ID: wishlistID, UserID: userID, Title: "Test"}
	mockRepo.On("GetByID", mock.Anything, wishlistID).Return(wishlist, nil)

	result, err := service.GetByID(context.Background(), wishlistID, userID)
	assert.NoError(t, err)
	assert.Equal(t, wishlist, result)
	mockRepo.AssertExpectations(t)
}

func TestWishlistService_GetByID_Forbidden(t *testing.T) {
	mockRepo := new(MockWishlistRepository)
	service := NewWishlistService(mockRepo)

	wishlistID := uuid.New()
	ownerID := uuid.New()
	otherUserID := uuid.New()
	wishlist := &models.Wishlist{ID: wishlistID, UserID: ownerID, Title: "Test"}
	mockRepo.On("GetByID", mock.Anything, wishlistID).Return(wishlist, nil)

	_, err := service.GetByID(context.Background(), wishlistID, otherUserID)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrForbidden))
}

func TestWishlistService_Update_Success(t *testing.T) {
	mockRepo := new(MockWishlistRepository)
	service := NewWishlistService(mockRepo)

	wishlistID := uuid.New()
	userID := uuid.New()
	wishlist := &models.Wishlist{ID: wishlistID, UserID: userID, Title: "Old", Description: "Old desc", EventDate: time.Now().AddDate(0, 1, 0)}
	mockRepo.On("GetByID", mock.Anything, wishlistID).Return(wishlist, nil)
	mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(w *models.Wishlist) bool {
		return w.Title == "New Title"
	})).Return(nil)

	err := service.Update(context.Background(), wishlistID, userID, "New Title", "New desc", wishlist.EventDate)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestWishlistService_Delete_Success(t *testing.T) {
	mockRepo := new(MockWishlistRepository)
	service := NewWishlistService(mockRepo)

	wishlistID := uuid.New()
	userID := uuid.New()
	wishlist := &models.Wishlist{ID: wishlistID, UserID: userID}
	mockRepo.On("GetByID", mock.Anything, wishlistID).Return(wishlist, nil)
	mockRepo.On("Delete", mock.Anything, wishlistID).Return(nil)

	err := service.Delete(context.Background(), wishlistID, userID)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
