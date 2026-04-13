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

// Mock репозитория
type mockWishlistRepo struct {
	mock.Mock
}

func (m *mockWishlistRepo) Create(ctx context.Context, userID int, title, description string, eventDate time.Time) (*models.Wishlist, error) {
	args := m.Called(ctx, userID, title, description, eventDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Wishlist), args.Error(1)
}

func (m *mockWishlistRepo) GetByID(ctx context.Context, id int) (*models.Wishlist, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Wishlist), args.Error(1)
}

func (m *mockWishlistRepo) GetByAccessToken(ctx context.Context, token uuid.UUID) (*models.Wishlist, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Wishlist), args.Error(1)
}

func (m *mockWishlistRepo) GetAllByUser(ctx context.Context, userID int) ([]models.Wishlist, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.Wishlist), args.Error(1)
}

func (m *mockWishlistRepo) Update(ctx context.Context, w *models.Wishlist) error {
	args := m.Called(ctx, w)
	return args.Error(0)
}

func (m *mockWishlistRepo) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestWishlistService_Create_Success(t *testing.T) {
	mockRepo := new(mockWishlistRepo)
	service := NewWishlistService(mockRepo)

	userID := 1
	title := "Birthday"
	desc := "Gifts"
	eventDate := time.Now().AddDate(0, 1, 0) // через месяц

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
	mockRepo := new(mockWishlistRepo)
	service := NewWishlistService(mockRepo)

	eventDate := time.Now().AddDate(0, 0, -1) // вчера
	_, err := service.Create(context.Background(), 1, "Title", "", eventDate)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "future")
}
