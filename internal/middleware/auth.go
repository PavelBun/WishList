// Package middleware provides HTTP middleware for the application.
package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"wishlist-api/internal/service"
)

type contextKey string

// UserIDKey is the context key for storing the authenticated user ID.
const UserIDKey contextKey = "userID"

// AuthMiddleware returns a middleware that validates JWT tokens and adds the user ID to the request context.
func AuthMiddleware(authService *service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
				return
			}
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Invalid Authorization header", http.StatusUnauthorized)
				return
			}
			token := parts[1]
			userID, err := authService.ValidateToken(token)
			if err != nil {
				http.Error(w, fmt.Sprintf("Invalid token: %v", err), http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
