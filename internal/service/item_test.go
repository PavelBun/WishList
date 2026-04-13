package service

import (
	"context"
	"testing"
	"time"
	"wishlist-api/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Полная реализация мока WishlistRepository
type mockWishlistRepoFull struct {
	mock.Mock
}

func (m *mockWishlistRepoFull) Create(ctx context.Context, userID int, title, description string, eventDate time.Time) (*models.Wishlist, error) {
	args := m.Called(ctx, userID, title, description, eventDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Wishlist), args.Error(1)
}

func (m *mockWishlistRepoFull) GetByID(ctx context.Context, id int) (*models.Wishlist, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Wishlist), args.Error(1)
}

func (m *mockWishlistRepoFull) GetByAccessToken(ctx context.Context, token uuid.UUID) (*models.Wishlist, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Wishlist), args.Error(1)
}

func (m *mockWishlistRepoFull) GetAllByUser(ctx context.Context, userID int) ([]models.Wishlist, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.Wishlist), args.Error(1)
}

func (m *mockWishlistRepoFull) Update(ctx context.Context, w *models.Wishlist) error {
	args := m.Called(ctx, w)
	return args.Error(0)
}

func (m *mockWishlistRepoFull) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Мок для ItemRepository (уже был)
type mockItemRepoFull struct {
	mock.Mock
}

func (m *mockItemRepoFull) Create(ctx context.Context, item *models.Item) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}
func (m *mockItemRepoFull) GetByID(ctx context.Context, id int) (*models.Item, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Item), args.Error(1)
}
func (m *mockItemRepoFull) GetAllByWishlistID(ctx context.Context, wishlistID int) ([]models.Item, error) {
	args := m.Called(ctx, wishlistID)
	return args.Get(0).([]models.Item), args.Error(1)
}
func (m *mockItemRepoFull) Update(ctx context.Context, item *models.Item) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}
func (m *mockItemRepoFull) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *mockItemRepoFull) BookItem(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestItemService_Create_Success(t *testing.T) {
	mockItemRepo := new(mockItemRepoFull)
	mockWishlistRepo := new(mockWishlistRepoFull)
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

func TestItemService_BookItem(t *testing.T) {
	mockItemRepo := new(mockItemRepoFull)
	mockWishlistRepo := new(mockWishlistRepoFull)
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
