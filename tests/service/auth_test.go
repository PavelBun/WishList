package service_test

import (
	"context"
	"testing"
	"time"

	"wishlist-api/internal/models"
	"wishlist-api/internal/service"
	"wishlist-api/pkg/hash"

	"github.com/stretchr/testify/assert"
)

type mockUserRepo struct {
	createFunc     func(ctx context.Context, email, passwordHash string) (*models.User, error)
	getByEmailFunc func(ctx context.Context, email string) (*models.User, error)
}

func (m *mockUserRepo) Create(ctx context.Context, email, passwordHash string) (*models.User, error) {
	return m.createFunc(ctx, email, passwordHash)
}
func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	return m.getByEmailFunc(ctx, email)
}

func TestAuthService_Register(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepo{
		getByEmailFunc: func(_ context.Context, _ string) (*models.User, error) {
			return nil, nil
		},
		createFunc: func(_ context.Context, email, _ string) (*models.User, error) {
			return &models.User{
				ID:        1,
				Email:     email,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}, nil
		},
	}
	authService := service.NewAuthService(repo, "test-secret")

	user, err := authService.Register(ctx, "test@example.com", "password123")
	assert.NoError(t, err)
	assert.Equal(t, "test@example.com", user.Email)
}

func TestAuthService_Login(t *testing.T) {
	ctx := context.Background()
	hashed, _ := hash.Password("password123")
	repo := &mockUserRepo{
		getByEmailFunc: func(_ context.Context, _ string) (*models.User, error) {
			return &models.User{
				ID:           1,
				Email:        "test@example.com",
				PasswordHash: hashed,
			}, nil
		},
	}
	authService := service.NewAuthService(repo, "test-secret")

	token, err := authService.Login(ctx, "test@example.com", "password123")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}
