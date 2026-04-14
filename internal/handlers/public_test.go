package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"wishlist-api/internal/models"
	"wishlist-api/internal/service"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

//nolint:funlen
func TestPublicHandler_GetWishlistByToken(t *testing.T) {
	mockWishlistSvc := new(MockWishlistService)
	mockItemSvc := new(MockItemService)
	handler := NewPublicHandler(mockWishlistSvc, mockItemSvc)

	token := uuid.New()
	wishlist := &models.Wishlist{
		ID:          uuid.New(),
		AccessToken: token,
		Title:       "Public Wishlist",
	}
	items := []models.Item{
		{ID: uuid.New(), WishlistID: wishlist.ID, Title: "Public Item"},
	}

	t.Run("success", func(t *testing.T) {
		mockWishlistSvc.On("GetByAccessToken", mock.Anything, token).Return(wishlist, nil).Once()
		mockItemSvc.On("GetAllPublicByWishlistID", mock.Anything, wishlist.ID).Return(items, nil).Once()

		req := newTestRequest(t, testRequest{
			method: http.MethodGet,
			path:   "/public/wishlists/" + token.String(),
			urlParams: map[string]string{
				"token": token.String(),
			},
		})
		w := httptest.NewRecorder()

		handler.GetWishlistByToken(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp models.Wishlist
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, wishlist.Title, resp.Title)
		assert.Len(t, resp.Items, 1)
	})

	t.Run("invalid token format", func(t *testing.T) {
		req := newTestRequest(t, testRequest{
			method: http.MethodGet,
			path:   "/public/wishlists/not-a-uuid",
			urlParams: map[string]string{
				"token": "not-a-uuid",
			},
		})
		w := httptest.NewRecorder()

		handler.GetWishlistByToken(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid token format")
	})

	t.Run("wishlist not found", func(t *testing.T) {
		mockWishlistSvc.On("GetByAccessToken", mock.Anything, token).Return(nil, service.ErrNotFound).Once()

		req := newTestRequest(t, testRequest{
			method: http.MethodGet,
			path:   "/public/wishlists/" + token.String(),
			urlParams: map[string]string{
				"token": token.String(),
			},
		})
		w := httptest.NewRecorder()

		handler.GetWishlistByToken(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestPublicHandler_BookItem(t *testing.T) {
	mockWishlistSvc := new(MockWishlistService)
	mockItemSvc := new(MockItemService)
	handler := NewPublicHandler(mockWishlistSvc, mockItemSvc)

	token := uuid.New()
	wishlistID := uuid.New()
	itemID := uuid.New()
	wishlist := &models.Wishlist{ID: wishlistID, AccessToken: token}

	t.Run("success", func(t *testing.T) {
		mockWishlistSvc.On("GetByAccessToken", mock.Anything, token).Return(wishlist, nil).Once()
		mockItemSvc.On("BookItem", mock.Anything, itemID, wishlistID).Return(nil).Once()

		req := newTestRequest(t, testRequest{
			method: http.MethodPost,
			path:   "/public/wishlists/" + token.String() + "/items/" + itemID.String() + "/book",
			urlParams: map[string]string{
				"token":   token.String(),
				"item_id": itemID.String(),
			},
		})
		w := httptest.NewRecorder()

		handler.BookItem(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("already booked", func(t *testing.T) {
		mockWishlistSvc.On("GetByAccessToken", mock.Anything, token).Return(wishlist, nil).Once()
		mockItemSvc.On("BookItem", mock.Anything, itemID, wishlistID).Return(service.ErrAlreadyBooked).Once()

		req := newTestRequest(t, testRequest{
			method: http.MethodPost,
			path:   "/public/wishlists/" + token.String() + "/items/" + itemID.String() + "/book",
			urlParams: map[string]string{
				"token":   token.String(),
				"item_id": itemID.String(),
			},
		})
		w := httptest.NewRecorder()

		handler.BookItem(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
	})
}
