package handlers

import (
	"context"
	"time"
	"wishlist-api/internal/models"
	"wishlist-api/internal/service"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockAuthService реализует service.AuthServiceInterface.
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(ctx context.Context, email, password string) (*models.User, error) {
	args := m.Called(ctx, email, password)
	if args.Get(0) == nil {
		return nil, args.Error(1) //nolint:wrapcheck
	}
	return args.Get(0).(*models.User), args.Error(1) //nolint:wrapcheck
}

func (m *MockAuthService) Login(ctx context.Context, email, password string) (string, error) {
	args := m.Called(ctx, email, password)
	return args.String(0), args.Error(1) //nolint:wrapcheck
}

func (m *MockAuthService) ValidateToken(tokenString string) (uuid.UUID, error) {
	args := m.Called(tokenString)
	return args.Get(0).(uuid.UUID), args.Error(1) //nolint:wrapcheck
}

// MockWishlistService реализует service.WishlistServiceInterface.
type MockWishlistService struct {
	mock.Mock
}

func (m *MockWishlistService) Create(ctx context.Context, userID uuid.UUID, title, description string, eventDate time.Time) (*models.Wishlist, error) {
	args := m.Called(ctx, userID, title, description, eventDate)
	if args.Get(0) == nil {
		return nil, args.Error(1) //nolint:wrapcheck
	}
	return args.Get(0).(*models.Wishlist), args.Error(1) //nolint:wrapcheck
}

func (m *MockWishlistService) GetByID(ctx context.Context, id, userID uuid.UUID) (*models.Wishlist, error) {
	args := m.Called(ctx, id, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1) //nolint:wrapcheck
	}
	return args.Get(0).(*models.Wishlist), args.Error(1) //nolint:wrapcheck
}

func (m *MockWishlistService) GetAllByUser(ctx context.Context, userID uuid.UUID) ([]models.Wishlist, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.Wishlist), args.Error(1) //nolint:wrapcheck
}

func (m *MockWishlistService) Update(ctx context.Context, id, userID uuid.UUID, title *string, description *string, eventDate *time.Time) error {
	args := m.Called(ctx, id, userID, title, description, eventDate)
	return args.Error(0) //nolint:wrapcheck
}

func (m *MockWishlistService) Delete(ctx context.Context, id, userID uuid.UUID) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0) //nolint:wrapcheck
}

func (m *MockWishlistService) GetByAccessToken(ctx context.Context, token uuid.UUID) (*models.Wishlist, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1) //nolint:wrapcheck
	}
	return args.Get(0).(*models.Wishlist), args.Error(1) //nolint:wrapcheck
}

// MockItemService реализует service.ItemServiceInterface.
type MockItemService struct {
	mock.Mock
}

func (m *MockItemService) Create(ctx context.Context, wishlistID, userID uuid.UUID, title, description, productLink string, priority models.Priority) (*models.Item, error) {
	args := m.Called(ctx, wishlistID, userID, title, description, productLink, priority)
	if args.Get(0) == nil {
		return nil, args.Error(1) //nolint:wrapcheck
	}
	return args.Get(0).(*models.Item), args.Error(1) //nolint:wrapcheck
}

func (m *MockItemService) GetByID(ctx context.Context, id, userID uuid.UUID) (*models.Item, error) {
	args := m.Called(ctx, id, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1) //nolint:wrapcheck
	}
	return args.Get(0).(*models.Item), args.Error(1) //nolint:wrapcheck
}

func (m *MockItemService) GetAllByWishlistID(ctx context.Context, wishlistID, userID uuid.UUID) ([]models.Item, error) {
	args := m.Called(ctx, wishlistID, userID)
	return args.Get(0).([]models.Item), args.Error(1) //nolint:wrapcheck
}

func (m *MockItemService) GetAllPublicByWishlistID(ctx context.Context, wishlistID uuid.UUID) ([]models.Item, error) {
	args := m.Called(ctx, wishlistID)
	return args.Get(0).([]models.Item), args.Error(1) //nolint:wrapcheck
}

func (m *MockItemService) Update(ctx context.Context, id, userID uuid.UUID, title *string, description *string, productLink *string, priority *models.Priority) error {
	args := m.Called(ctx, id, userID, title, description, productLink, priority)
	return args.Error(0) //nolint:wrapcheck
}

func (m *MockItemService) Delete(ctx context.Context, id, userID uuid.UUID) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0) //nolint:wrapcheck
}

func (m *MockItemService) BookItem(ctx context.Context, id, wishlistID uuid.UUID) error {
	args := m.Called(ctx, id, wishlistID)
	return args.Error(0) //nolint:wrapcheck
}

// Compilation checks for interface implementation
var _ service.AuthServiceInterface = (*MockAuthService)(nil)
var _ service.WishlistServiceInterface = (*MockWishlistService)(nil)
var _ service.ItemServiceInterface = (*MockItemService)(nil)
