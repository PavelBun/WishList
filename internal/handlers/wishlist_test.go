package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"wishlist-api/internal/dto"
	"wishlist-api/internal/models"
	"wishlist-api/internal/service"

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

	t.Run("success", func(t *testing.T) {
		expectedWishlist := &models.Wishlist{
			ID:     uuid.New(),
			UserID: userID,
			Title:  validReq.Title,
		}
		mockWishlistSvc.On("Create", mock.Anything, userID, validReq.Title, validReq.Description, mock.AnythingOfType("time.Time")).
			Return(expectedWishlist, nil).Once()

		req := newTestRequest(t, testRequest{
			method: http.MethodPost,
			path:   "/wishlists",
			body:   validReq,
			userID: userID,
		})
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
		req := newTestRequest(t, testRequest{
			method: http.MethodPost,
			path:   "/wishlists",
			body:   validReq,
		})
		w := httptest.NewRecorder()

		handler.Create(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "unauthorized")
	})

	t.Run("invalid event date format", func(t *testing.T) {
		invalidReq := validReq
		invalidReq.EventDate = "not-a-date"

		req := newTestRequest(t, testRequest{
			method: http.MethodPost,
			path:   "/wishlists",
			body:   invalidReq,
			userID: userID,
		})
		w := httptest.NewRecorder()

		handler.Create(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "YYYY-MM-DD")
	})

	t.Run("service error", func(t *testing.T) {
		mockWishlistSvc.On("Create", mock.Anything, userID, validReq.Title, validReq.Description, mock.AnythingOfType("time.Time")).
			Return(nil, assert.AnError).Once()

		req := newTestRequest(t, testRequest{
			method: http.MethodPost,
			path:   "/wishlists",
			body:   validReq,
			userID: userID,
		})
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

		req := newTestRequest(t, testRequest{
			method: http.MethodGet,
			path:   "/wishlists",
			userID: userID,
		})
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
		req := newTestRequest(t, testRequest{
			method: http.MethodGet,
			path:   "/wishlists",
		})
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

		req := newTestRequest(t, testRequest{
			method: http.MethodGet,
			path:   "/wishlists/" + wishlistID.String(),
			userID: userID,
			urlParams: map[string]string{
				"id": wishlistID.String(),
			},
		})
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

	t.Run("not found", func(t *testing.T) {
		mockWishlistSvc.On("GetByID", mock.Anything, wishlistID, userID).Return(nil, service.ErrNotFound).Once()

		req := newTestRequest(t, testRequest{
			method: http.MethodGet,
			path:   "/wishlists/" + wishlistID.String(),
			userID: userID,
			urlParams: map[string]string{
				"id": wishlistID.String(),
			},
		})
		w := httptest.NewRecorder()

		handler.GetByID(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestWishlistHandler_Update(t *testing.T) {
	mockWishlistSvc := new(MockWishlistService)
	mockItemSvc := new(MockItemService)
	handler := NewWishlistHandler(mockWishlistSvc, mockItemSvc)

	userID := uuid.New()
	wishlistID := uuid.New()
	updateReq := dto.UpdateWishlistRequest{
		Title:       strPtr("New Title"),
		Description: strPtr("New Desc"),
	}
	expectedWishlist := &models.Wishlist{
		ID:     wishlistID,
		UserID: userID,
		Title:  "New Title",
	}

	t.Run("success", func(t *testing.T) {
		mockWishlistSvc.On("Update", mock.Anything, wishlistID, userID,
			updateReq.Title, updateReq.Description, mock.Anything).
			Return(nil).Once()
		mockWishlistSvc.On("GetByID", mock.Anything, wishlistID, userID).
			Return(expectedWishlist, nil).Once()

		req := newTestRequest(t, testRequest{
			method: http.MethodPut,
			path:   "/wishlists/" + wishlistID.String(),
			body:   updateReq,
			userID: userID,
			urlParams: map[string]string{
				"id": wishlistID.String(),
			},
		})
		w := httptest.NewRecorder()

		handler.Update(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp models.Wishlist
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, expectedWishlist.Title, resp.Title)
		mockWishlistSvc.AssertExpectations(t)
	})

	t.Run("unauthorized", func(t *testing.T) {
		req := newTestRequest(t, testRequest{
			method: http.MethodPut,
			path:   "/wishlists/" + wishlistID.String(),
			body:   updateReq,
			urlParams: map[string]string{
				"id": wishlistID.String(),
			},
		})
		w := httptest.NewRecorder()

		handler.Update(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
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

		req := newTestRequest(t, testRequest{
			method: http.MethodDelete,
			path:   "/wishlists/" + wishlistID.String(),
			userID: userID,
			urlParams: map[string]string{
				"id": wishlistID.String(),
			},
		})
		w := httptest.NewRecorder()

		handler.Delete(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
		mockWishlistSvc.AssertExpectations(t)
	})

	t.Run("unauthorized", func(t *testing.T) {
		req := newTestRequest(t, testRequest{
			method: http.MethodDelete,
			path:   "/wishlists/" + wishlistID.String(),
			urlParams: map[string]string{
				"id": wishlistID.String(),
			},
		})
		w := httptest.NewRecorder()

		handler.Delete(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

// strPtr returns a pointer to a string.
func strPtr(s string) *string {
	return &s
}
