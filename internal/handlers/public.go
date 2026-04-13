package handlers

import (
	"net/http"
	"strconv"
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
// @Param item_id path int true "Item ID"
// @Success 200 {string} string "booked"
// @Failure 409 {string} string "item already booked"
// @Router /public/wishlists/{token}/items/{item_id}/book [post]
func (h *PublicHandler) BookItem(w http.ResponseWriter, r *http.Request) {
	tokenStr := chi.URLParam(r, "token")
	token, err := uuid.Parse(tokenStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid token format")
		return
	}
	itemIDStr := chi.URLParam(r, "item_id")
	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid item ID")
		return
	}

	wishlist, err := h.wishlistService.GetByAccessToken(r.Context(), token)
	if err != nil {
		writeSafeError(w, r, err)
		return
	}

	err = h.itemService.BookItem(r.Context(), itemID, wishlist.ID)
	if err != nil {
		writeSafeError(w, r, err)
		return
	}
	writeJSONSuccess(w, map[string]string{"status": "booked"})
}
