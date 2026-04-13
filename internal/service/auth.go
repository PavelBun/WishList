// Package service contains business logic.
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
	"github.com/google/uuid"
)

// Domain errors.
var (
	ErrNotFound           = errors.New("resource not found")
	ErrForbidden          = errors.New("access denied")
	ErrAlreadyBooked      = errors.New("item already booked")
	ErrInvalidInput       = errors.New("invalid input")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserExists         = errors.New("user already exists")
)

// AuthService handles authentication logic.
type AuthService struct {
	userRepo  repository.UserRepository
	jwtSecret string
	jwtExpiry time.Duration
}

// NewAuthService creates a new AuthService.
func NewAuthService(userRepo repository.UserRepository, jwtSecret string, expiryHours int) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
		jwtExpiry: time.Duration(expiryHours) * time.Hour,
	}
}

// Register creates a new user account.
func (s *AuthService) Register(ctx context.Context, email, password string) (*models.User, error) {
	existing, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}
	if existing != nil {
		return nil, ErrUserExists
	}
	passwordHash, err := hash.Password(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	user, err := s.userRepo.Create(ctx, email, passwordHash)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return user, nil
}

// Login authenticates a user and returns a JWT token.
func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil || user == nil {
		return "", ErrInvalidCredentials
	}
	if !hash.Check(password, user.PasswordHash) {
		return "", ErrInvalidCredentials
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.String(),
		"email":   user.Email,
		"exp":     time.Now().Add(s.jwtExpiry).Unix(),
	})
	signed, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return signed, nil
}

// ValidateToken parses and validates a JWT token, returning the user ID.
func (s *AuthService) ValidateToken(tokenString string) (uuid.UUID, error) {
	token, err := jwt.Parse(tokenString, func(_ *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})
	if err != nil || !token.Valid {
		return uuid.Nil, ErrInvalidCredentials
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, ErrInvalidCredentials
	}
	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return uuid.Nil, ErrInvalidCredentials
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, ErrInvalidCredentials
	}
	return userID, nil
}
