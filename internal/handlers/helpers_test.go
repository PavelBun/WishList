package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"wishlist-api/internal/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// testRequest holds data for building a test HTTP request.
type testRequest struct {
	method    string
	path      string
	body      interface{}
	userID    uuid.UUID
	urlParams map[string]string
}

// newTestRequest creates an HTTP request with context and chi URL parameters.
func newTestRequest(tb testingTB, req testRequest) *http.Request {
	var bodyReader *bytes.Reader
	if req.body != nil {
		bodyBytes, err := json.Marshal(req.body)
		if err != nil {
			tb.Fatalf("failed to marshal request body: %v", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	} else {
		bodyReader = bytes.NewReader(nil)
	}

	r := httptest.NewRequest(req.method, req.path, bodyReader)
	r.Header.Set("Content-Type", "application/json")

	// Добавляем userID в контекст, используя правильный ключ
	if req.userID != uuid.Nil {
		ctx := context.WithValue(r.Context(), middleware.UserIDKey, req.userID)
		r = r.WithContext(ctx)
	}

	// Добавляем chi URL параметры
	if len(req.urlParams) > 0 {
		rctx := chi.NewRouteContext()
		for key, val := range req.urlParams {
			rctx.URLParams.Add(key, val)
		}
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	}

	return r
}

// testingTB is a minimal interface covering *testing.T and *testing.B.
type testingTB interface {
	Fatalf(format string, args ...interface{})
}
