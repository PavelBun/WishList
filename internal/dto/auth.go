// Package dto contains Data Transfer Objects for API requests and responses.
package dto

// RegisterRequest represents the request body for user registration.
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// LoginRequest represents the request body for user login.
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse represents the response body for a successful login.
type LoginResponse struct {
	Token string `json:"token"`
}
