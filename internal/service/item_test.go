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

// setupItemTest creates mock repositories, service instance, and common test IDs.
func setupItemTest(t *testing.T) (*MockItemRepository, *MockWishlistRepository, *ItemService, uuid.UUID, uuid.UUID) {
	t.Helper()
	mockItemRepo := new(MockItemRepository)
	mockWishlistRepo := new(MockWishlistRepository)
	svc := NewItemService(mockItemRepo, mockWishlistRepo)
	userID := uuid.New()
	wishlistID := uuid.New()
	return mockItemRepo, mockWishlistRepo, svc, userID, wishlistID
}

// expectWishlistGet configures the wishlist mock to return a valid wishlist owned by userID.
func expectWishlistGet(mockWishlistRepo *MockWishlistRepository, wishlistID, userID uuid.UUID) {
	mockWishlistRepo.On("GetByID", mock.Anything, wishlistID).
		Return(&models.Wishlist{ID: wishlistID, UserID: userID}, nil)
}

// expectWishlistNotFound configures the wishlist mock to return ErrNotFound.
func expectWishlistNotFound(mockWishlistRepo *MockWishlistRepository, wishlistID uuid.UUID) {
	mockWishlistRepo.On("GetByID", mock.Anything, wishlistID).
		Return(nil, repository.ErrNotFound)
}

// ==================== Create tests ====================

func TestItemService_Create_Success(t *testing.T) {
	mockItemRepo, mockWishlistRepo, svc, userID, wishlistID := setupItemTest(t)
	title := "PS5"
	desc := "Game console"
	link := "https://example.com"
	priority := models.PriorityMust

	expectWishlistGet(mockWishlistRepo, wishlistID, userID)
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
	_, _, svc, _, wishlistID := setupItemTest(t)
	_, err := svc.Create(context.Background(), wishlistID, uuid.New(), "title", "", "", models.Priority(0))
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrInvalidInput))
}

func TestItemService_Create_WishlistNotFound(t *testing.T) {
	_, mockWishlistRepo, svc, userID, wishlistID := setupItemTest(t)
	expectWishlistNotFound(mockWishlistRepo, wishlistID)

	_, err := svc.Create(context.Background(), wishlistID, userID, "title", "", "", models.PriorityLow)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrNotFound))
}

func TestItemService_Create_Forbidden(t *testing.T) {
	_, mockWishlistRepo, svc, _, wishlistID := setupItemTest(t)
	otherUserID := uuid.New()
	// Wishlist belongs to a different user
	mockWishlistRepo.On("GetByID", mock.Anything, wishlistID).
		Return(&models.Wishlist{ID: wishlistID, UserID: uuid.New()}, nil)

	_, err := svc.Create(context.Background(), wishlistID, otherUserID, "title", "", "", models.PriorityLow)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrForbidden))
}

// ==================== GetByID tests ====================

func TestItemService_GetByID_Success(t *testing.T) {
	mockItemRepo, mockWishlistRepo, svc, userID, wishlistID := setupItemTest(t)
	itemID := uuid.New()
	item := &models.Item{ID: itemID, WishlistID: wishlistID}

	mockItemRepo.On("GetByID", mock.Anything, itemID).Return(item, nil)
	expectWishlistGet(mockWishlistRepo, wishlistID, userID)

	result, err := svc.GetByID(context.Background(), itemID, userID)
	assert.NoError(t, err)
	assert.Equal(t, item, result)
}

func TestItemService_GetByID_NotFound(t *testing.T) {
	mockItemRepo, _, svc, userID, _ := setupItemTest(t)
	itemID := uuid.New()
	mockItemRepo.On("GetByID", mock.Anything, itemID).Return(nil, repository.ErrNotFound)

	_, err := svc.GetByID(context.Background(), itemID, userID)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrNotFound))
}

func TestItemService_GetByID_Forbidden(t *testing.T) {
	mockItemRepo, mockWishlistRepo, svc, _, wishlistID := setupItemTest(t)
	itemID := uuid.New()
	otherUserID := uuid.New()
	item := &models.Item{ID: itemID, WishlistID: wishlistID}

	mockItemRepo.On("GetByID", mock.Anything, itemID).Return(item, nil)
	// Wishlist belongs to another user
	mockWishlistRepo.On("GetByID", mock.Anything, wishlistID).
		Return(&models.Wishlist{ID: wishlistID, UserID: uuid.New()}, nil)

	_, err := svc.GetByID(context.Background(), itemID, otherUserID)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrForbidden))
}

// ==================== Update tests ====================

func TestItemService_Update_Success(t *testing.T) {
	mockItemRepo, mockWishlistRepo, svc, userID, wishlistID := setupItemTest(t)
	itemID := uuid.New()
	item := &models.Item{
		ID:         itemID,
		WishlistID: wishlistID,
		Title:      "Old",
		Priority:   models.PriorityMedium,
	}
	expectWishlistGet(mockWishlistRepo, wishlistID, userID)
	mockItemRepo.On("GetByID", mock.Anything, itemID).Return(item, nil)

	newTitle := "New Title"
	newPriority := models.PriorityMust
	mockItemRepo.On("Update", mock.Anything, mock.MatchedBy(func(i *models.Item) bool {
		return i.Title == newTitle && i.Priority == newPriority
	})).Return(nil)

	err := svc.Update(context.Background(), itemID, userID, &newTitle, nil, nil, &newPriority)
	assert.NoError(t, err)
	mockItemRepo.AssertExpectations(t)
}

func TestItemService_Update_NotFound(t *testing.T) {
	mockItemRepo, _, svc, userID, _ := setupItemTest(t)
	itemID := uuid.New()
	mockItemRepo.On("GetByID", mock.Anything, itemID).Return(nil, repository.ErrNotFound)

	newTitle := "New Title"
	err := svc.Update(context.Background(), itemID, userID, &newTitle, nil, nil, nil)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrNotFound))
}

func TestItemService_Update_Forbidden(t *testing.T) {
	mockItemRepo, mockWishlistRepo, svc, _, wishlistID := setupItemTest(t)
	itemID := uuid.New()
	otherUserID := uuid.New()
	item := &models.Item{ID: itemID, WishlistID: wishlistID}
	mockItemRepo.On("GetByID", mock.Anything, itemID).Return(item, nil)
	// Wishlist belongs to another user
	mockWishlistRepo.On("GetByID", mock.Anything, wishlistID).
		Return(&models.Wishlist{ID: wishlistID, UserID: uuid.New()}, nil)

	newTitle := "New Title"
	err := svc.Update(context.Background(), itemID, otherUserID, &newTitle, nil, nil, nil)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrForbidden))
}

func TestItemService_Update_InvalidPriority(t *testing.T) {
	mockItemRepo, mockWishlistRepo, svc, userID, wishlistID := setupItemTest(t)
	itemID := uuid.New()
	item := &models.Item{ID: itemID, WishlistID: wishlistID}
	expectWishlistGet(mockWishlistRepo, wishlistID, userID)
	mockItemRepo.On("GetByID", mock.Anything, itemID).Return(item, nil)

	invalidPriority := models.Priority(0)
	err := svc.Update(context.Background(), itemID, userID, nil, nil, nil, &invalidPriority)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrInvalidInput))
}

// ==================== Delete tests ====================

func TestItemService_Delete_Success(t *testing.T) {
	mockItemRepo, mockWishlistRepo, svc, userID, wishlistID := setupItemTest(t)
	itemID := uuid.New()
	item := &models.Item{ID: itemID, WishlistID: wishlistID}
	expectWishlistGet(mockWishlistRepo, wishlistID, userID)
	mockItemRepo.On("GetByID", mock.Anything, itemID).Return(item, nil)
	mockItemRepo.On("Delete", mock.Anything, itemID).Return(nil)

	err := svc.Delete(context.Background(), itemID, userID)
	assert.NoError(t, err)
}

func TestItemService_Delete_NotFound(t *testing.T) {
	mockItemRepo, _, svc, userID, _ := setupItemTest(t)
	itemID := uuid.New()
	mockItemRepo.On("GetByID", mock.Anything, itemID).Return(nil, repository.ErrNotFound)

	err := svc.Delete(context.Background(), itemID, userID)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrNotFound))
}

func TestItemService_Delete_Forbidden(t *testing.T) {
	mockItemRepo, mockWishlistRepo, svc, _, wishlistID := setupItemTest(t)
	itemID := uuid.New()
	otherUserID := uuid.New()
	item := &models.Item{ID: itemID, WishlistID: wishlistID}
	mockItemRepo.On("GetByID", mock.Anything, itemID).Return(item, nil)
	// Wishlist belongs to another user
	mockWishlistRepo.On("GetByID", mock.Anything, wishlistID).
		Return(&models.Wishlist{ID: wishlistID, UserID: uuid.New()}, nil)

	err := svc.Delete(context.Background(), itemID, otherUserID)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrForbidden))
}

// ==================== BookItem tests ====================

func TestItemService_BookItem_Success(t *testing.T) {
	mockItemRepo, _, svc, _, wishlistID := setupItemTest(t)
	itemID := uuid.New()
	item := &models.Item{ID: itemID, WishlistID: wishlistID}
	mockItemRepo.On("GetByID", mock.Anything, itemID).Return(item, nil)
	mockItemRepo.On("BookItem", mock.Anything, itemID).Return(nil)

	err := svc.BookItem(context.Background(), itemID, wishlistID)
	assert.NoError(t, err)
}

func TestItemService_BookItem_AlreadyBooked(t *testing.T) {
	mockItemRepo, _, svc, _, wishlistID := setupItemTest(t)
	itemID := uuid.New()
	item := &models.Item{ID: itemID, WishlistID: wishlistID}
	mockItemRepo.On("GetByID", mock.Anything, itemID).Return(item, nil)
	mockItemRepo.On("BookItem", mock.Anything, itemID).Return(repository.ErrAlreadyBooked)

	err := svc.BookItem(context.Background(), itemID, wishlistID)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrAlreadyBooked))
}

func TestItemService_BookItem_NotFound(t *testing.T) {
	mockItemRepo, _, svc, _, wishlistID := setupItemTest(t)
	itemID := uuid.New()
	mockItemRepo.On("GetByID", mock.Anything, itemID).Return(nil, repository.ErrNotFound)

	err := svc.BookItem(context.Background(), itemID, wishlistID)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrNotFound))
}

func TestItemService_BookItem_WrongWishlist(t *testing.T) {
	mockItemRepo, _, svc, _, _ := setupItemTest(t)
	itemID := uuid.New()
	wishlistID := uuid.New()
	wrongWishlistID := uuid.New()
	item := &models.Item{ID: itemID, WishlistID: wishlistID}
	mockItemRepo.On("GetByID", mock.Anything, itemID).Return(item, nil)

	err := svc.BookItem(context.Background(), itemID, wrongWishlistID)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrForbidden))
}

// ==================== GetAllByWishlistID tests ====================

func TestItemService_GetAllByWishlistID_Success(t *testing.T) {
	mockItemRepo, mockWishlistRepo, svc, userID, wishlistID := setupItemTest(t)
	expectWishlistGet(mockWishlistRepo, wishlistID, userID)
	expectedItems := []models.Item{{ID: uuid.New(), Title: "Item1"}}
	mockItemRepo.On("GetAllByWishlistID", mock.Anything, wishlistID).Return(expectedItems, nil)

	items, err := svc.GetAllByWishlistID(context.Background(), wishlistID, userID)
	assert.NoError(t, err)
	assert.Equal(t, expectedItems, items)
}

func TestItemService_GetAllByWishlistID_Forbidden(t *testing.T) {
	_, mockWishlistRepo, svc, _, wishlistID := setupItemTest(t)
	otherUserID := uuid.New()
	mockWishlistRepo.On("GetByID", mock.Anything, wishlistID).
		Return(&models.Wishlist{ID: wishlistID, UserID: uuid.New()}, nil)

	_, err := svc.GetAllByWishlistID(context.Background(), wishlistID, otherUserID)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrForbidden))
}
