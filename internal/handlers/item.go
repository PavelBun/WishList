package handlers

import (
	"net/http"
	"wishlist-api/internal/dto"
	"wishlist-api/internal/middleware"
	"wishlist-api/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// ItemHandler handles item-related endpoints.
type ItemHandler struct {
	itemService *service.ItemService
}

// NewItemHandler creates a new ItemHandler instance.
func NewItemHandler(itemService *service.ItemService) *ItemHandler {
	return &ItemHandler{itemService: itemService}
}

// Create godoc
// @Summary Create item in wishlist
// @Tags items
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param wishlist_id path string true "Wishlist ID (UUID)"
// @Param request body dto.CreateItemRequest true "Item data"
// @Success 201 {object} models.Item
// @Router /wishlists/{wishlist_id}/items [post]
func (h *ItemHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	wishlistIDStr := chi.URLParam(r, "wishlist_id")
	wishlistID, err := uuid.Parse(wishlistIDStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid wishlist ID")
		return
	}
	var req dto.CreateItemRequest
	if err := decodeAndValidate(r, &req); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	item, err := h.itemService.Create(r.Context(), wishlistID, userID, req.Title, req.Description, req.ProductLink, req.Priority)
	if err != nil {
		writeSafeError(w, r, err)
		return
	}
	writeJSONCreated(w, item)
}

// GetAll godoc
// @Summary Get all items of a wishlist
// @Tags items
// @Security BearerAuth
// @Produce json
// @Param wishlist_id path string true "Wishlist ID (UUID)"
// @Success 200 {array} models.Item
// @Router /wishlists/{wishlist_id}/items [get]
func (h *ItemHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	wishlistIDStr := chi.URLParam(r, "wishlist_id")
	wishlistID, err := uuid.Parse(wishlistIDStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid wishlist ID")
		return
	}
	items, err := h.itemService.GetAllByWishlistID(r.Context(), wishlistID, userID)
	if err != nil {
		writeSafeError(w, r, err)
		return
	}
	writeJSONSuccess(w, items)
}

// GetByID godoc
// @Summary Get a single item by ID
// @Tags items
// @Security BearerAuth
// @Produce json
// @Param id path string true "Item ID (UUID)"
// @Success 200 {object} models.Item
// @Failure 400 {string} string "Invalid item ID"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Access denied"
// @Failure 404 {string} string "Item not found"
// @Router /items/{id} [get]
func (h *ItemHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid item ID")
		return
	}
	item, err := h.itemService.GetByID(r.Context(), id, userID)
	if err != nil {
		writeSafeError(w, r, err)
		return
	}
	writeJSONSuccess(w, item)
}

// Update godoc
// @Summary Update item
// @Tags items
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Item ID (UUID)"
// @Param request body dto.UpdateItemRequest true "Updated fields"
// @Success 200 {object} models.Item
// @Failure 400 {string} string "Invalid input"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Access denied"
// @Failure 404 {string} string "Item not found"
// @Router /items/{id} [put]
func (h *ItemHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid item ID")
		return
	}
	var req dto.UpdateItemRequest
	if err := decodeAndValidate(r, &req); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.itemService.Update(r.Context(), id, userID, req.Title, req.Description, req.ProductLink, req.Priority); err != nil {
		writeSafeError(w, r, err)
		return
	}

	updated, err := h.itemService.GetByID(r.Context(), id, userID)
	if err != nil {
		// This should not happen after successful update, but handle gracefully
		writeJSONSuccess(w, map[string]string{"status": "updated"})
		return
	}
	writeJSONSuccess(w, updated)
}

// Delete godoc
// @Summary Delete item
// @Tags items
// @Security BearerAuth
// @Param id path string true "Item ID (UUID)"
// @Success 204 "No Content"
// @Failure 400 {string} string "Invalid item ID"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Access denied"
// @Failure 404 {string} string "Item not found"
// @Router /items/{id} [delete]
func (h *ItemHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid item ID")
		return
	}
	if err := h.itemService.Delete(r.Context(), id, userID); err != nil {
		writeSafeError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
