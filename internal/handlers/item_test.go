package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"wishlist-api/internal/dto"
	"wishlist-api/internal/middleware"
	"wishlist-api/internal/models"
	"wishlist-api/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

//nolint:funlen
func TestItemHandler_Create(t *testing.T) {
	mockItemSvc := new(MockItemService)
	handler := NewItemHandler(mockItemSvc)

	userID := uuid.New()
	wishlistID := uuid.New()
	validReq := dto.CreateItemRequest{
		Title:       "Laptop",
		Description: "MacBook Pro",
		ProductLink: "https://example.com",
		Priority:    models.PriorityMust,
	}
	validBody, _ := json.Marshal(validReq)

	t.Run("success", func(t *testing.T) {
		expectedItem := &models.Item{
			ID:         uuid.New(),
			WishlistID: wishlistID,
			Title:      validReq.Title,
		}
		mockItemSvc.On("Create", mock.Anything, wishlistID, userID,
			validReq.Title, validReq.Description, validReq.ProductLink, validReq.Priority).
			Return(expectedItem, nil).Once()

		req := httptest.NewRequest(http.MethodPost, "/wishlists/"+wishlistID.String()+"/items", bytes.NewReader(validBody))
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("wishlist_id", wishlistID.String())
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		handler.Create(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		var resp models.Item
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, expectedItem.Title, resp.Title)
		mockItemSvc.AssertExpectations(t)
	})

	t.Run("unauthorized", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/wishlists/"+wishlistID.String()+"/items", bytes.NewReader(validBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler.Create(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "unauthorized")
	})

	t.Run("invalid wishlist ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/wishlists/invalid/items", bytes.NewReader(validBody))
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("wishlist_id", "invalid")
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		handler.Create(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid wishlist ID")
	})
}

func TestItemHandler_GetAll(t *testing.T) {
	mockItemSvc := new(MockItemService)
	handler := NewItemHandler(mockItemSvc)

	userID := uuid.New()
	wishlistID := uuid.New()
	expectedItems := []models.Item{
		{ID: uuid.New(), WishlistID: wishlistID, Title: "Item1"},
		{ID: uuid.New(), WishlistID: wishlistID, Title: "Item2"},
	}

	t.Run("success", func(t *testing.T) {
		mockItemSvc.On("GetAllByWishlistID", mock.Anything, wishlistID, userID).
			Return(expectedItems, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/wishlists/"+wishlistID.String()+"/items", nil)
		ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("wishlist_id", wishlistID.String())
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		handler.GetAll(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp []models.Item
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Len(t, resp, 2)
		mockItemSvc.AssertExpectations(t)
	})

	t.Run("unauthorized", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/wishlists/"+wishlistID.String()+"/items", nil)
		w := httptest.NewRecorder()
		handler.GetAll(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestItemHandler_GetByID(t *testing.T) {
	mockItemSvc := new(MockItemService)
	handler := NewItemHandler(mockItemSvc)

	userID := uuid.New()
	itemID := uuid.New()
	expectedItem := &models.Item{ID: itemID, Title: "Test Item"}

	t.Run("success", func(t *testing.T) {
		mockItemSvc.On("GetByID", mock.Anything, itemID, userID).Return(expectedItem, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/items/"+itemID.String(), nil)
		ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", itemID.String())
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		handler.GetByID(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp models.Item
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, expectedItem.Title, resp.Title)
	})

	t.Run("not found", func(t *testing.T) {
		mockItemSvc.On("GetByID", mock.Anything, itemID, userID).Return(nil, service.ErrNotFound).Once()

		req := httptest.NewRequest(http.MethodGet, "/items/"+itemID.String(), nil)
		ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", itemID.String())
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		handler.GetByID(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestItemHandler_Delete(t *testing.T) {
	mockItemSvc := new(MockItemService)
	handler := NewItemHandler(mockItemSvc)

	userID := uuid.New()
	itemID := uuid.New()

	t.Run("success", func(t *testing.T) {
		mockItemSvc.On("Delete", mock.Anything, itemID, userID).Return(nil).Once()

		req := httptest.NewRequest(http.MethodDelete, "/items/"+itemID.String(), nil)
		ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", itemID.String())
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		handler.Delete(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})
}
