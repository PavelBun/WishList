package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"wishlist-api/internal/dto"
	"wishlist-api/internal/middleware"
	"wishlist-api/internal/service"
	"wishlist-api/internal/validator"

	"github.com/go-chi/chi/v5"
)

// WishlistHandler handles wishlist endpoints.
type WishlistHandler struct {
	wishlistService *service.WishlistService
}

// NewWishlistHandler creates a new WishlistHandler.
func NewWishlistHandler(wishlistService *service.WishlistService) *WishlistHandler {
	return &WishlistHandler{wishlistService: wishlistService}
}

// Create creates a new wishlist.
// @Summary Create new wishlist
// @Tags wishlists
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateWishlistRequest true "Wishlist data"
// @Success 201 {object} models.Wishlist
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Router /wishlists [post]
func (h *WishlistHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)
	var req dto.CreateWishlistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if err := validator.Validate(req); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	wishlist, err := h.wishlistService.Create(r.Context(), userID, req.Title, req.Description, req.EventDate)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSONCreated(w, wishlist)
}

// GetAll returns all wishlists for the authenticated user.
// @Summary Get all user's wishlists
// @Tags wishlists
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.Wishlist
// @Failure 401 {string} string
// @Router /wishlists [get]
func (h *WishlistHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)
	wishlists, err := h.wishlistService.GetAllByUser(r.Context(), userID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSONSuccess(w, wishlists)
}

// GetByID returns a single wishlist by ID.
// @Summary Get wishlist by ID
// @Tags wishlists
// @Security BearerAuth
// @Produce json
// @Param id path int true "Wishlist ID"
// @Success 200 {object} models.Wishlist
// @Failure 404 {string} string
// @Router /wishlists/{id} [get]
func (h *WishlistHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid ID")
		return
	}
	wishlist, err := h.wishlistService.GetByID(r.Context(), id, userID)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSONSuccess(w, wishlist)
}

// Update updates an existing wishlist.
// @Summary Update wishlist
// @Tags wishlists
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Wishlist ID"
// @Param request body dto.UpdateWishlistRequest true "Updated fields"
// @Success 200 {string} string "OK"
// @Failure 400,404,401 {string} string
// @Router /wishlists/{id} [put]
func (h *WishlistHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid ID")
		return
	}
	var req dto.UpdateWishlistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if err := validator.Validate(req); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.wishlistService.Update(r.Context(), id, userID, req.Title, req.Description, req.EventDate); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSONSuccess(w, map[string]string{"status": "updated"})
}

// Delete deletes a wishlist.
// @Summary Delete wishlist
// @Tags wishlists
// @Security BearerAuth
// @Param id path int true "Wishlist ID"
// @Success 200 {string} string "OK"
// @Failure 400,404,401 {string} string
// @Router /wishlists/{id} [delete]
func (h *WishlistHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid ID")
		return
	}
	if err := h.wishlistService.Delete(r.Context(), id, userID); err != nil {
		writeJSONError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSONSuccess(w, map[string]string{"status": "deleted"})
}
