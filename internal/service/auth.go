package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"wishlist-api/internal/models"
	"wishlist-api/internal/repository"
	"wishlist-api/pkg/hash"

	"github.com/golang-jwt/jwt/v5"
)

// AuthService handles user authentication and token management.
type AuthService struct {
	userRepo  repository.UserRepository
	jwtSecret string
}

// NewAuthService creates a new AuthService.
func NewAuthService(userRepo repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

// Register creates a new user account.
func (s *AuthService) Register(ctx context.Context, email, password string) (*models.User, error) {
	existing, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}
	if existing != nil {
		return nil, errors.New("user already exists")
	}
	passwordHash, err := hash.Password(password)
	if err != nil {
		return nil, err
	}
	return s.userRepo.Create(ctx, email, passwordHash)
}

// Login authenticates a user and returns a JWT token.
func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil || user == nil {
		return "", errors.New("invalid email or password")
	}
	if !hash.Check(password, user.PasswordHash) {
		return "", errors.New("invalid email or password")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})
	return token.SignedString([]byte(s.jwtSecret))
}

// ValidateToken validates a JWT token and returns the user ID.
func (s *AuthService) ValidateToken(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(_ *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})
	if err != nil || !token.Valid {
		return 0, errors.New("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid claims")
	}
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("invalid user_id in token")
	}
	return int(userID), nil
}
