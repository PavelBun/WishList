// Package handlers contains HTTP request handlers.
package handlers

import (
	"net/http"
	"wishlist-api/internal/dto"
	"wishlist-api/internal/service"
)

// AuthHandler handles authentication endpoints.
type AuthHandler struct {
	authService service.AuthServiceInterface
}

// NewAuthHandler creates a new AuthHandler instance.
func NewAuthHandler(authService service.AuthServiceInterface) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register godoc
// @Summary Register new user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "User data"
// @Success 201 {object} models.User
// @Failure 400 {string} string
// @Failure 409 {string} string
// @Router /auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := decodeAndValidate(r, &req); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	user, err := h.authService.Register(r.Context(), req.Email, req.Password)
	if err != nil {
		writeSafeError(w, r, err)
		return
	}
	writeJSONCreated(w, user)
}

// Login godoc
// @Summary Login user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Credentials"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Router /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := decodeAndValidate(r, &req); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	token, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		writeSafeError(w, r, err)
		return
	}
	writeJSONSuccess(w, dto.LoginResponse{Token: token})
}
