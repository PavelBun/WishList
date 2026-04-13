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
)

// Domain errors.
var (
	// ErrNotFound is returned when a requested resource does not exist.
	ErrNotFound = errors.New("resource not found")
	// ErrForbidden is returned when access to a resource is denied.
	ErrForbidden = errors.New("access denied")
	// ErrAlreadyBooked is returned when trying to book an already booked item.
	ErrAlreadyBooked = errors.New("item already booked")
	// ErrInvalidInput is returned when input validation fails.
	ErrInvalidInput = errors.New("invalid input")
	// ErrInvalidCredentials is returned when email or password is incorrect.
	ErrInvalidCredentials = errors.New("invalid email or password")
	// ErrUserExists is returned when trying to register an existing user.
	ErrUserExists = errors.New("user already exists")
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
		return nil, err
	}
	return s.userRepo.Create(ctx, email, passwordHash)
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
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(s.jwtExpiry).Unix(),
	})
	return token.SignedString([]byte(s.jwtSecret))
}

// ValidateToken parses and validates a JWT token, returning the user ID.
func (s *AuthService) ValidateToken(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(_ *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})
	if err != nil || !token.Valid {
		return 0, ErrInvalidCredentials
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, ErrInvalidCredentials
	}
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, ErrInvalidCredentials
	}
	return int(userID), nil
}
