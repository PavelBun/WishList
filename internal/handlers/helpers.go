// Package handlers provides HTTP request handlers for the Wishlist API.
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"wishlist-api/internal/validator"
)

// writeJSON sends a JSON response with the given status code.
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// writeJSONSuccess sends a 200 OK JSON response.
func writeJSONSuccess(w http.ResponseWriter, data any) {
	writeJSON(w, http.StatusOK, data)
}

// writeJSONCreated sends a 201 Created JSON response.
func writeJSONCreated(w http.ResponseWriter, data any) {
	writeJSON(w, http.StatusCreated, data)
}

// writeJSONError sends a JSON error response with the given status code and message.
func writeJSONError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

// decodeAndValidate decodes JSON from request body and validates the struct.
func decodeAndValidate(r *http.Request, dst interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return fmt.Errorf("invalid request body")
	}
	if err := validator.Validate(dst); err != nil {
		return err
	}
	return nil
}
