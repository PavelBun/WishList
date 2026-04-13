package handlers

import (
	"net/http"
	"strconv"
	"wishlist-api/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// PublicHandler handles public endpoints that do not require authentication.
type PublicHandler struct {
	wishlistService *service.WishlistService
	itemService     *service.ItemService
}

// NewPublicHandler creates a new PublicHandler.
func NewPublicHandler(wishlistService *service.WishlistService, itemService *service.ItemService) *PublicHandler {
	return &PublicHandler{wishlistService: wishlistService, itemService: itemService}
}

// GetWishlistByToken returns a public wishlist by its access token.
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
	if err != nil || wishlist == nil {
		writeJSONError(w, http.StatusNotFound, "Wishlist not found")
		return
	}
	items, err := h.itemService.GetAllPublicByWishlistID(r.Context(), wishlist.ID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to load wishlist items")
		return
	}
	wishlist.Items = items
	writeJSONSuccess(w, wishlist)
}

// BookItem books an item in a public wishlist.
// @Summary Book an item in public wishlist
// @Tags public
// @Param token path string true "Access token (UUID)"
// @Param item_id path int true "Item ID"
// @Success 200 {string} string "booked"
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
	if err != nil || wishlist == nil {
		writeJSONError(w, http.StatusNotFound, "Wishlist not found")
		return
	}

	err = h.itemService.BookItem(r.Context(), itemID, wishlist.ID)
	if err != nil {
		writeJSONError(w, http.StatusConflict, err.Error())
		return
	}
	writeJSONSuccess(w, map[string]string{"status": "booked"})
}
