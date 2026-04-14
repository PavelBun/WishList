package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"wishlist-api/internal/dto"
	"wishlist-api/internal/models"
	"wishlist-api/internal/service"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

//nolint:funlen
func TestAuthHandler_Register(t *testing.T) {
	mockSvc := new(MockAuthService)
	handler := NewAuthHandler(mockSvc)

	validReq := dto.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	validBody, _ := json.Marshal(validReq)

	t.Run("success", func(t *testing.T) {
		expectedUser := &models.User{
			ID:    uuid.New(),
			Email: validReq.Email,
		}
		mockSvc.On("Register", mock.Anything, validReq.Email, validReq.Password).
			Return(expectedUser, nil).Once()

		req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(validBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Register(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		var respUser models.User
		err := json.Unmarshal(w.Body.Bytes(), &respUser)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser.Email, respUser.Email)
		mockSvc.AssertExpectations(t)
	})

	t.Run("invalid request body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader([]byte("invalid")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Register(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid request body")
	})

	t.Run("validation error", func(t *testing.T) {
		invalidReq := dto.RegisterRequest{
			Email:    "not-an-email",
			Password: "short",
		}
		body, _ := json.Marshal(invalidReq)
		req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Register(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "validation failed")
	})

	t.Run("user already exists", func(t *testing.T) {
		mockSvc.On("Register", mock.Anything, validReq.Email, validReq.Password).
			Return(nil, service.ErrUserExists).Once()

		req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(validBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Register(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
		assert.Contains(t, w.Body.String(), service.ErrUserExists.Error())
		mockSvc.AssertExpectations(t)
	})
}

func TestAuthHandler_Login(t *testing.T) {
	mockSvc := new(MockAuthService)
	handler := NewAuthHandler(mockSvc)

	validReq := dto.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	validBody, _ := json.Marshal(validReq)

	t.Run("success", func(t *testing.T) {
		mockSvc.On("Login", mock.Anything, validReq.Email, validReq.Password).
			Return("jwt-token", nil).Once()

		req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(validBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Login(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp dto.LoginResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "jwt-token", resp.Token)
		mockSvc.AssertExpectations(t)
	})

	t.Run("invalid credentials", func(t *testing.T) {
		mockSvc.On("Login", mock.Anything, validReq.Email, validReq.Password).
			Return("", service.ErrInvalidCredentials).Once()

		req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(validBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Login(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), service.ErrInvalidCredentials.Error())
		mockSvc.AssertExpectations(t)
	})
}
