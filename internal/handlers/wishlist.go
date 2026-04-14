package handlers

import (
	"net/http"
	"time"
	"wishlist-api/internal/dto"
	"wishlist-api/internal/middleware"
	"wishlist-api/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// WishlistHandler handles wishlist-related endpoints.
type WishlistHandler struct {
	wishlistService *service.WishlistService
	itemService     *service.ItemService
}

// NewWishlistHandler creates a new WishlistHandler instance.
func NewWishlistHandler(wishlistService *service.WishlistService, itemService *service.ItemService) *WishlistHandler {
	return &WishlistHandler{
		wishlistService: wishlistService,
		itemService:     itemService,
	}
}

// Create godoc
// @Summary Create new wishlist
// @Tags wishlists
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateWishlistRequest true "Wishlist data"
// @Success 201 {object} models.Wishlist
// @Router /wishlists [post]
func (h *WishlistHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	var req dto.CreateWishlistRequest
	if err := decodeAndValidate(r, &req); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	eventDate, err := time.Parse("2006-01-02", req.EventDate)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "event_date must be in YYYY-MM-DD format")
		return
	}

	wishlist, err := h.wishlistService.Create(r.Context(), userID, req.Title, req.Description, eventDate)
	if err != nil {
		writeSafeError(w, r, err)
		return
	}
	writeJSONCreated(w, wishlist)
}

// GetAll godoc
// @Summary Get all user's wishlists
// @Tags wishlists
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.Wishlist
// @Router /wishlists [get]
func (h *WishlistHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	wishlists, err := h.wishlistService.GetAllByUser(r.Context(), userID)
	if err != nil {
		writeSafeError(w, r, err)
		return
	}
	writeJSONSuccess(w, wishlists)
}

// GetByID godoc
// @Summary Get wishlist by ID with items
// @Tags wishlists
// @Security BearerAuth
// @Produce json
// @Param id path string true "Wishlist ID (UUID)"
// @Success 200 {object} models.Wishlist
// @Router /wishlists/{id} [get]
func (h *WishlistHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid ID")
		return
	}
	wishlist, err := h.wishlistService.GetByID(r.Context(), id, userID)
	if err != nil {
		writeSafeError(w, r, err)
		return
	}
	// Загружаем предметы
	items, err := h.itemService.GetAllByWishlistID(r.Context(), id, userID)
	if err != nil {
		writeSafeError(w, r, err)
		return
	}
	wishlist.Items = items
	writeJSONSuccess(w, wishlist)
}

// Update godoc
// @Summary Update wishlist
// @Tags wishlists
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Wishlist ID (UUID)"
// @Param request body dto.UpdateWishlistRequest true "Updated fields"
// @Success 200 {object} models.Wishlist
// @Router /wishlists/{id} [put]
func (h *WishlistHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid ID")
		return
	}
	var req dto.UpdateWishlistRequest
	if err := decodeAndValidate(r, &req); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	var eventDatePtr *time.Time
	if req.EventDate != nil {
		parsed, err := time.Parse("2006-01-02", *req.EventDate)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "event_date must be in YYYY-MM-DD format")
			return
		}
		eventDatePtr = &parsed
	}

	if err := h.wishlistService.Update(r.Context(), id, userID, req.Title, req.Description, eventDatePtr); err != nil {
		writeSafeError(w, r, err)
		return
	}

	updated, err := h.wishlistService.GetByID(r.Context(), id, userID)
	if err != nil {
		writeJSONSuccess(w, map[string]string{"status": "updated"})
		return
	}
	writeJSONSuccess(w, updated)
}

// Delete godoc
// @Summary Delete wishlist
// @Tags wishlists
// @Security BearerAuth
// @Param id path string true "Wishlist ID (UUID)"
// @Success 204 "No Content"
// @Router /wishlists/{id} [delete]
func (h *WishlistHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid ID")
		return
	}
	if err := h.wishlistService.Delete(r.Context(), id, userID); err != nil {
		writeSafeError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
