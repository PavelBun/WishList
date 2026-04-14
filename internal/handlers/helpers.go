package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"wishlist-api/internal/middleware"
	"wishlist-api/internal/validator"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// writeJSON writes any JSON response.
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("failed to encode JSON response", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// writeJSONSuccess writes a 200 OK response with data.
func writeJSONSuccess(w http.ResponseWriter, data any) {
	writeJSON(w, http.StatusOK, data)
}

// writeJSONCreated writes a 201 Created response with data.
func writeJSONCreated(w http.ResponseWriter, data any) {
	writeJSON(w, http.StatusCreated, data)
}

// decodeAndValidate decodes JSON from request body and validates the struct.
func decodeAndValidate(r *http.Request, dst interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return fmt.Errorf("invalid request body: %w", err)
	}
	if err := validator.Validate(dst); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}
	return nil
}

// getUserID extracts the authenticated user ID from request context.
func getUserID(r *http.Request) (uuid.UUID, error) {
	id, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("user ID not found in context")
	}
	return id, nil
}

// parseUUIDParam extracts and parses a UUID path parameter.
func parseUUIDParam(r *http.Request, name string) (uuid.UUID, error) {
	param := chi.URLParam(r, name)
	if param == "" {
		return uuid.Nil, fmt.Errorf("missing path parameter %s", name)
	}
	id, err := uuid.Parse(param)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid UUID for %s: %w", name, err)
	}
	return id, nil
}
