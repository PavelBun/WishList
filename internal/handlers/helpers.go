package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"wishlist-api/internal/validator"
)

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("failed to encode JSON response", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func writeJSONSuccess(w http.ResponseWriter, data any) {
	writeJSON(w, http.StatusOK, data)
}

func writeJSONCreated(w http.ResponseWriter, data any) {
	writeJSON(w, http.StatusCreated, data)
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	slog.Error("API error", "status", status, "message", message)
	writeJSON(w, status, map[string]string{"error": message})
}

func decodeAndValidate(r *http.Request, dst interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return fmt.Errorf("invalid request body")
	}
	if err := validator.Validate(dst); err != nil {
		return err
	}
	return nil
}
