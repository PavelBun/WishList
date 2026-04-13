package service

import (
	"context"
	"time"
	"wishlist-api/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockWishlistRepository is a mock implementation of repository.WishlistRepository.
type MockWishlistRepository struct {
	mock.Mock
}

// Create mocks the Create method.
func (m *MockWishlistRepository) Create(ctx context.Context, userID uuid.UUID, title, description string, eventDate time.Time) (*models.Wishlist, error) {
	args := m.Called(ctx, userID, title, description, eventDate)
	if args.Get(0) == nil {
		return nil, args.Error(1) //nolint:wrapcheck
	}
	return args.Get(0).(*models.Wishlist), args.Error(1) //nolint:wrapcheck
}

// GetByID mocks GetByID.
func (m *MockWishlistRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Wishlist, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1) //nolint:wrapcheck
	}
	return args.Get(0).(*models.Wishlist), args.Error(1) //nolint:wrapcheck
}

// GetByAccessToken mocks GetByAccessToken.
func (m *MockWishlistRepository) GetByAccessToken(ctx context.Context, token uuid.UUID) (*models.Wishlist, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1) //nolint:wrapcheck
	}
	return args.Get(0).(*models.Wishlist), args.Error(1) //nolint:wrapcheck
}

// GetAllByUser mocks GetAllByUser.
func (m *MockWishlistRepository) GetAllByUser(ctx context.Context, userID uuid.UUID) ([]models.Wishlist, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.Wishlist), args.Error(1) //nolint:wrapcheck
}

// Update mocks Update.
func (m *MockWishlistRepository) Update(ctx context.Context, w *models.Wishlist) error {
	args := m.Called(ctx, w)
	return args.Error(0) //nolint:wrapcheck
}

// Delete mocks Delete.
func (m *MockWishlistRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0) //nolint:wrapcheck
}

// MockItemRepository is a mock implementation of repository.ItemRepository.
type MockItemRepository struct {
	mock.Mock
}

// Create mocks Create.
func (m *MockItemRepository) Create(ctx context.Context, item *models.Item) error {
	args := m.Called(ctx, item)
	return args.Error(0) //nolint:wrapcheck
}

// GetByID mocks GetByID.
func (m *MockItemRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Item, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1) //nolint:wrapcheck
	}
	return args.Get(0).(*models.Item), args.Error(1) //nolint:wrapcheck
}

// GetAllByWishlistID mocks GetAllByWishlistID.
func (m *MockItemRepository) GetAllByWishlistID(ctx context.Context, wishlistID uuid.UUID) ([]models.Item, error) {
	args := m.Called(ctx, wishlistID)
	return args.Get(0).([]models.Item), args.Error(1) //nolint:wrapcheck
}

// Update mocks Update.
func (m *MockItemRepository) Update(ctx context.Context, item *models.Item) error {
	args := m.Called(ctx, item)
	return args.Error(0) //nolint:wrapcheck
}

// Delete mocks Delete.
func (m *MockItemRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0) //nolint:wrapcheck
}

// BookItem mocks BookItem.
func (m *MockItemRepository) BookItem(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0) //nolint:wrapcheck
}
