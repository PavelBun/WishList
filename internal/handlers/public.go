package handlers

import (
	"net/http"
	"wishlist-api/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// PublicHandler handles public endpoints (no auth).
type PublicHandler struct {
	wishlistService *service.WishlistService
	itemService     *service.ItemService
}

// NewPublicHandler creates a new PublicHandler instance.
func NewPublicHandler(wishlistService *service.WishlistService, itemService *service.ItemService) *PublicHandler {
	return &PublicHandler{wishlistService: wishlistService, itemService: itemService}
}

// GetWishlistByToken godoc
// @Summary Get public wishlist by access token
// @Tags public
// @Produce json
// @Param token path string true "Access token (UUID)"
// @Success 200 {object} models.Wishlist
// @Failure 400 {string} string "Invalid token format"
// @Failure 404 {string} string "Wishlist not found"
// @Router /public/wishlists/{token} [get]
func (h *PublicHandler) GetWishlistByToken(w http.ResponseWriter, r *http.Request) {
	tokenStr := chi.URLParam(r, "token")
	token, err := uuid.Parse(tokenStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid token format")
		return
	}
	wishlist, err := h.wishlistService.GetByAccessToken(r.Context(), token)
	if err != nil {
		writeSafeError(w, r, err)
		return
	}
	items, err := h.itemService.GetAllPublicByWishlistID(r.Context(), wishlist.ID)
	if err != nil {
		writeSafeError(w, r, err)
		return
	}
	wishlist.Items = items
	writeJSONSuccess(w, wishlist)
}

// BookItem godoc
// @Summary Book an item in public wishlist
// @Tags public
// @Param token path string true "Access token (UUID)"
// @Param item_id path string true "Item ID (UUID)"
// @Success 204 "No Content"
// @Failure 400 {string} string "Invalid token or item ID"
// @Failure 404 {string} string "Wishlist or item not found"
// @Failure 409 {string} string "Item already booked"
// @Router /public/wishlists/{token}/items/{item_id}/book [post]
func (h *PublicHandler) BookItem(w http.ResponseWriter, r *http.Request) {
	tokenStr := chi.URLParam(r, "token")
	token, err := uuid.Parse(tokenStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid token format")
		return
	}
	itemIDStr := chi.URLParam(r, "item_id")
	itemID, err := uuid.Parse(itemIDStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid item ID")
		return
	}
	wishlist, err := h.wishlistService.GetByAccessToken(r.Context(), token)
	if err != nil {
		writeSafeError(w, r, err)
		return
	}
	if err := h.itemService.BookItem(r.Context(), itemID, wishlist.ID); err != nil {
		writeSafeError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
