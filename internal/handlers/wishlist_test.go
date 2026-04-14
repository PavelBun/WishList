package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"wishlist-api/internal/dto"
	"wishlist-api/internal/middleware"
	"wishlist-api/internal/models"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

//nolint:funlen
func TestWishlistHandler_Create(t *testing.T) {
	mockWishlistSvc := new(MockWishlistService)
	mockItemSvc := new(MockItemService)
	handler := NewWishlistHandler(mockWishlistSvc, mockItemSvc)

	userID := uuid.New()
	validReq := dto.CreateWishlistRequest{
		Title:       "Birthday",
		Description: "Gifts",
		EventDate:   time.Now().AddDate(0, 1, 0).Format("2006-01-02"),
	}
	validBody, _ := json.Marshal(validReq)

	t.Run("success", func(t *testing.T) {
		expectedWishlist := &models.Wishlist{
			ID:     uuid.New(),
			UserID: userID,
			Title:  validReq.Title,
		}
		mockWishlistSvc.On("Create", mock.Anything, userID, validReq.Title, validReq.Description, mock.AnythingOfType("time.Time")).
			Return(expectedWishlist, nil).Once()

		req := httptest.NewRequest(http.MethodPost, "/wishlists", bytes.NewReader(validBody))
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		handler.Create(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		var resp models.Wishlist
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, expectedWishlist.Title, resp.Title)
		mockWishlistSvc.AssertExpectations(t)
	})

	t.Run("unauthorized - missing userID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/wishlists", bytes.NewReader(validBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Create(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "unauthorized")
	})

	t.Run("invalid event date format", func(t *testing.T) {
		invalidReq := validReq
		invalidReq.EventDate = "not-a-date"
		body, _ := json.Marshal(invalidReq)

		req := httptest.NewRequest(http.MethodPost, "/wishlists", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		handler.Create(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "YYYY-MM-DD")
	})

	t.Run("service error", func(t *testing.T) {
		mockWishlistSvc.On("Create", mock.Anything, userID, validReq.Title, validReq.Description, mock.AnythingOfType("time.Time")).
			Return(nil, assert.AnError).Once()

		req := httptest.NewRequest(http.MethodPost, "/wishlists", bytes.NewReader(validBody))
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		handler.Create(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockWishlistSvc.AssertExpectations(t)
	})
}

func TestWishlistHandler_GetAll(t *testing.T) {
	mockWishlistSvc := new(MockWishlistService)
	mockItemSvc := new(MockItemService)
	handler := NewWishlistHandler(mockWishlistSvc, mockItemSvc)

	userID := uuid.New()
	expected := []models.Wishlist{
		{ID: uuid.New(), UserID: userID, Title: "List 1"},
		{ID: uuid.New(), UserID: userID, Title: "List 2"},
	}

	t.Run("success", func(t *testing.T) {
		mockWishlistSvc.On("GetAllByUser", mock.Anything, userID).Return(expected, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/wishlists", nil)
		ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		handler.GetAll(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp []models.Wishlist
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Len(t, resp, 2)
		mockWishlistSvc.AssertExpectations(t)
	})

	t.Run("unauthorized", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/wishlists", nil)
		w := httptest.NewRecorder()
		handler.GetAll(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "unauthorized")
	})
}

func TestWishlistHandler_GetByID(t *testing.T) {
	mockWishlistSvc := new(MockWishlistService)
	mockItemSvc := new(MockItemService)
	handler := NewWishlistHandler(mockWishlistSvc, mockItemSvc)

	userID := uuid.New()
	wishlistID := uuid.New()
	expectedWishlist := &models.Wishlist{ID: wishlistID, UserID: userID, Title: "Test"}
	expectedItems := []models.Item{{ID: uuid.New(), Title: "Item1"}}

	t.Run("success", func(t *testing.T) {
		mockWishlistSvc.On("GetByID", mock.Anything, wishlistID, userID).Return(expectedWishlist, nil).Once()
		mockItemSvc.On("GetAllByWishlistID", mock.Anything, wishlistID, userID).Return(expectedItems, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/wishlists/"+wishlistID.String(), nil)
		ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID)
		req = req.WithContext(ctx)
		// chi URL param routing
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", wishlistID.String())
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		w := httptest.NewRecorder()

		handler.GetByID(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp models.Wishlist
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, expectedWishlist.Title, resp.Title)
		assert.Len(t, resp.Items, 1)
		mockWishlistSvc.AssertExpectations(t)
		mockItemSvc.AssertExpectations(t)
	})
}

func TestWishlistHandler_Delete(t *testing.T) {
	mockWishlistSvc := new(MockWishlistService)
	mockItemSvc := new(MockItemService)
	handler := NewWishlistHandler(mockWishlistSvc, mockItemSvc)

	userID := uuid.New()
	wishlistID := uuid.New()

	t.Run("success", func(t *testing.T) {
		mockWishlistSvc.On("Delete", mock.Anything, wishlistID, userID).Return(nil).Once()

		req := httptest.NewRequest(http.MethodDelete, "/wishlists/"+wishlistID.String(), nil)
		ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID)
		req = req.WithContext(ctx)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", wishlistID.String())
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		w := httptest.NewRecorder()

		handler.Delete(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
		mockWishlistSvc.AssertExpectations(t)
	})
}
