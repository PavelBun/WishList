package handlers

import (
	"net/http"
	"strconv"
	"wishlist-api/internal/dto"
	"wishlist-api/internal/middleware"
	"wishlist-api/internal/service"

	"github.com/go-chi/chi/v5"
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
// @Param wishlist_id path int true "Wishlist ID"
// @Param request body dto.CreateItemRequest true "Item data"
// @Success 201 {object} models.Item
// @Router /wishlists/{wishlist_id}/items [post]
func (h *ItemHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)
	wishlistIDStr := chi.URLParam(r, "wishlist_id")
	wishlistID, err := strconv.Atoi(wishlistIDStr)
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
// @Param wishlist_id path int true "Wishlist ID"
// @Success 200 {array} models.Item
// @Router /wishlists/{wishlist_id}/items [get]
func (h *ItemHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)
	wishlistIDStr := chi.URLParam(r, "wishlist_id")
	wishlistID, err := strconv.Atoi(wishlistIDStr)
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

// Update godoc
// @Summary Update item
// @Tags items
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Item ID"
// @Param request body dto.UpdateItemRequest true "Updated fields"
// @Success 200 {string} string "OK"
// @Router /items/{id} [put]
func (h *ItemHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
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
	writeJSONSuccess(w, map[string]string{"status": "updated"})
}

// Delete godoc
// @Summary Delete item
// @Tags items
// @Security BearerAuth
// @Param id path int true "Item ID"
// @Success 200 {string} string "OK"
// @Router /items/{id} [delete]
func (h *ItemHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid item ID")
		return
	}
	if err := h.itemService.Delete(r.Context(), id, userID); err != nil {
		writeSafeError(w, r, err)
		return
	}
	writeJSONSuccess(w, map[string]string{"status": "deleted"})
}
