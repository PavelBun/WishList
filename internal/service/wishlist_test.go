package service

import (
	"context"
	"errors"
	"testing"
	"time"
	"wishlist-api/internal/models"
	"wishlist-api/internal/repository"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// setupWishlistTest creates mock repository, service instance, and common IDs.
func setupWishlistTest(t *testing.T) (*MockWishlistRepository, *WishlistService, uuid.UUID, uuid.UUID) {
	t.Helper()
	mockRepo := new(MockWishlistRepository)
	svc := NewWishlistService(mockRepo)
	userID := uuid.New()
	wishlistID := uuid.New()
	return mockRepo, svc, userID, wishlistID
}

// ==================== Create tests ====================

func TestWishlistService_Create_Success(t *testing.T) {
	mockRepo, svc, userID, _ := setupWishlistTest(t)
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
	mockRepo.On("Create", mock.Anything, userID, title, desc, eventDate).Return(expectedWishlist, nil)

	wishlist, err := svc.Create(context.Background(), userID, title, desc, eventDate)
	assert.NoError(t, err)
	assert.Equal(t, expectedWishlist, wishlist)
}

func TestWishlistService_Create_PastDate(t *testing.T) {
	_, svc, userID, _ := setupWishlistTest(t)
	eventDate := time.Now().AddDate(0, 0, -1)
	_, err := svc.Create(context.Background(), userID, "Title", "", eventDate)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrInvalidInput))
}

// ==================== GetByID tests ====================

func TestWishlistService_GetByID_Success(t *testing.T) {
	mockRepo, svc, userID, wishlistID := setupWishlistTest(t)
	wishlist := &models.Wishlist{ID: wishlistID, UserID: userID, Title: "Test"}
	mockRepo.On("GetByID", mock.Anything, wishlistID).Return(wishlist, nil)

	result, err := svc.GetByID(context.Background(), wishlistID, userID)
	assert.NoError(t, err)
	assert.Equal(t, wishlist, result)
}

func TestWishlistService_GetByID_Forbidden(t *testing.T) {
	mockRepo, svc, _, wishlistID := setupWishlistTest(t)
	otherUserID := uuid.New()
	wishlist := &models.Wishlist{ID: wishlistID, UserID: uuid.New(), Title: "Test"}
	mockRepo.On("GetByID", mock.Anything, wishlistID).Return(wishlist, nil)

	_, err := svc.GetByID(context.Background(), wishlistID, otherUserID)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrForbidden))
}

func TestWishlistService_GetByID_NotFound(t *testing.T) {
	mockRepo, svc, userID, wishlistID := setupWishlistTest(t)
	mockRepo.On("GetByID", mock.Anything, wishlistID).Return(nil, repository.ErrNotFound)

	_, err := svc.GetByID(context.Background(), wishlistID, userID)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrNotFound))
}

// ==================== Update tests ====================

func TestWishlistService_Update_Success(t *testing.T) {
	mockRepo, svc, userID, wishlistID := setupWishlistTest(t)
	wishlist := &models.Wishlist{
		ID:          wishlistID,
		UserID:      userID,
		Title:       "Old",
		Description: "Old desc",
		EventDate:   time.Now().AddDate(0, 1, 0),
	}
	mockRepo.On("GetByID", mock.Anything, wishlistID).Return(wishlist, nil)
	mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(w *models.Wishlist) bool {
		return w.Title == "New Title" && w.Description == "New desc"
	})).Return(nil)

	newTitle := "New Title"
	newDesc := "New desc"
	err := svc.Update(context.Background(), wishlistID, userID, &newTitle, &newDesc, &wishlist.EventDate)
	assert.NoError(t, err)
}

func TestWishlistService_Update_NotFound(t *testing.T) {
	mockRepo, svc, userID, wishlistID := setupWishlistTest(t)
	mockRepo.On("GetByID", mock.Anything, wishlistID).Return(nil, repository.ErrNotFound)

	newTitle := "title"
	newDesc := "desc"
	newDate := time.Now().AddDate(0, 1, 0)
	err := svc.Update(context.Background(), wishlistID, userID, &newTitle, &newDesc, &newDate)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrNotFound))
}

func TestWishlistService_Update_Forbidden(t *testing.T) {
	mockRepo, svc, _, wishlistID := setupWishlistTest(t)
	otherUserID := uuid.New()
	wishlist := &models.Wishlist{ID: wishlistID, UserID: uuid.New(), Title: "Old"}
	mockRepo.On("GetByID", mock.Anything, wishlistID).Return(wishlist, nil)

	newTitle := "New Title"
	err := svc.Update(context.Background(), wishlistID, otherUserID, &newTitle, nil, nil)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrForbidden))
}

func TestWishlistService_Update_InvalidDate(t *testing.T) {
	mockRepo, svc, userID, wishlistID := setupWishlistTest(t)
	wishlist := &models.Wishlist{ID: wishlistID, UserID: userID, EventDate: time.Now().AddDate(0, 1, 0)}
	mockRepo.On("GetByID", mock.Anything, wishlistID).Return(wishlist, nil)

	pastDate := time.Now().AddDate(0, 0, -1)
	err := svc.Update(context.Background(), wishlistID, userID, nil, nil, &pastDate)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrInvalidInput))
}

// ==================== Delete tests ====================

func TestWishlistService_Delete_Success(t *testing.T) {
	mockRepo, svc, userID, wishlistID := setupWishlistTest(t)
	wishlist := &models.Wishlist{ID: wishlistID, UserID: userID}
	mockRepo.On("GetByID", mock.Anything, wishlistID).Return(wishlist, nil)
	mockRepo.On("Delete", mock.Anything, wishlistID).Return(nil)

	err := svc.Delete(context.Background(), wishlistID, userID)
	assert.NoError(t, err)
}

func TestWishlistService_Delete_NotFound(t *testing.T) {
	mockRepo, svc, userID, wishlistID := setupWishlistTest(t)
	mockRepo.On("GetByID", mock.Anything, wishlistID).Return(nil, repository.ErrNotFound)

	err := svc.Delete(context.Background(), wishlistID, userID)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrNotFound))
}

func TestWishlistService_Delete_Forbidden(t *testing.T) {
	mockRepo, svc, _, wishlistID := setupWishlistTest(t)
	otherUserID := uuid.New()
	wishlist := &models.Wishlist{ID: wishlistID, UserID: uuid.New()}
	mockRepo.On("GetByID", mock.Anything, wishlistID).Return(wishlist, nil)

	err := svc.Delete(context.Background(), wishlistID, otherUserID)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrForbidden))
}

// ==================== GetAllByUser tests ====================

func TestWishlistService_GetAllByUser_Success(t *testing.T) {
	mockRepo, svc, userID, _ := setupWishlistTest(t)
	expected := []models.Wishlist{{ID: uuid.New(), UserID: userID}}
	mockRepo.On("GetAllByUser", mock.Anything, userID).Return(expected, nil)

	result, err := svc.GetAllByUser(context.Background(), userID)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

// ==================== GetByAccessToken tests ====================

func TestWishlistService_GetByAccessToken_Success(t *testing.T) {
	mockRepo, svc, _, _ := setupWishlistTest(t)
	token := uuid.New()
	wishlist := &models.Wishlist{ID: uuid.New(), AccessToken: token}
	mockRepo.On("GetByAccessToken", mock.Anything, token).Return(wishlist, nil)

	result, err := svc.GetByAccessToken(context.Background(), token)
	assert.NoError(t, err)
	assert.Equal(t, wishlist, result)
}

func TestWishlistService_GetByAccessToken_NotFound(t *testing.T) {
	mockRepo, svc, _, _ := setupWishlistTest(t)
	token := uuid.New()
	mockRepo.On("GetByAccessToken", mock.Anything, token).Return(nil, repository.ErrNotFound)

	_, err := svc.GetByAccessToken(context.Background(), token)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrNotFound))
}
