package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"wishlist-api/internal/middleware"
	"wishlist-api/internal/repository"
	"wishlist-api/internal/service"
)

// statusFromError maps domain/repository errors to HTTP status codes.
func statusFromError(err error) int {
	switch {
	case errors.Is(err, service.ErrNotFound),
		errors.Is(err, repository.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, service.ErrForbidden):
		return http.StatusForbidden
	case errors.Is(err, service.ErrAlreadyBooked):
		return http.StatusConflict
	case errors.Is(err, service.ErrInvalidInput):
		return http.StatusBadRequest
	case errors.Is(err, service.ErrInvalidCredentials):
		return http.StatusUnauthorized
	case errors.Is(err, service.ErrUserExists):
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

// writeSafeError logs the error with request ID and sends a safe message to the client.
func writeSafeError(w http.ResponseWriter, r *http.Request, err error) {
	status := statusFromError(err)
	requestID := middleware.GetRequestID(r.Context())
	message := err.Error()

	if status == http.StatusInternalServerError {
		slog.Error("internal server error",
			"path", r.URL.Path,
			"request_id", requestID,
			"error", err,
		)
		message = "internal server error"
	} else {
		slog.Warn("request error",
			"path", r.URL.Path,
			"status", status,
			"request_id", requestID,
			"error", err,
		)
	}
	writeJSONError(w, status, message)
}

// writeJSONError writes a JSON error response.
func writeJSONError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
