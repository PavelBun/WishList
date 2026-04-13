package service

import (
	"context"
	"time"
	"wishlist-api/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockWishlistRepository struct {
	mock.Mock
}

func (m *MockWishlistRepository) Create(ctx context.Context, userID int, title, description string, eventDate time.Time) (*models.Wishlist, error) {
	args := m.Called(ctx, userID, title, description, eventDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Wishlist), args.Error(1)
}

func (m *MockWishlistRepository) GetByID(ctx context.Context, id int) (*models.Wishlist, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Wishlist), args.Error(1)
}

func (m *MockWishlistRepository) GetByAccessToken(ctx context.Context, token uuid.UUID) (*models.Wishlist, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Wishlist), args.Error(1)
}

func (m *MockWishlistRepository) GetAllByUser(ctx context.Context, userID int) ([]models.Wishlist, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.Wishlist), args.Error(1)
}

func (m *MockWishlistRepository) Update(ctx context.Context, w *models.Wishlist) error {
	args := m.Called(ctx, w)
	return args.Error(0)
}

func (m *MockWishlistRepository) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockItemRepository struct {
	mock.Mock
}

func (m *MockItemRepository) Create(ctx context.Context, item *models.Item) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockItemRepository) GetByID(ctx context.Context, id int) (*models.Item, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Item), args.Error(1)
}

func (m *MockItemRepository) GetAllByWishlistID(ctx context.Context, wishlistID int) ([]models.Item, error) {
	args := m.Called(ctx, wishlistID)
	return args.Get(0).([]models.Item), args.Error(1)
}

func (m *MockItemRepository) Update(ctx context.Context, item *models.Item) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockItemRepository) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockItemRepository) BookItem(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
