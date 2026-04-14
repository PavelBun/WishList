package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"wishlist-api/internal/validator"
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
