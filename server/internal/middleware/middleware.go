// Package middleware - handles auth check
package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	auth "tanoserver/internal/auth"
)

var ErrUnauthorized = errors.New("unauthorized")

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authStr := r.Header.Get("Authorization")
		if authStr == "" {
			http.Error(
				w,
				"missing Authorization header",
				http.StatusBadRequest,
			)
			return
		}

		if !strings.HasPrefix(authStr, "Bearer ") {
			http.Error(
				w,
				"invalid Authorization header",
				http.StatusBadRequest,
			)
			return
		}

		tokenString := strings.TrimPrefix(authStr, "Bearer ")
		userID, err := auth.VerifyToken(tokenString)
		if err != nil {
			http.Error(
				w,
				"invalid token",
				http.StatusUnauthorized,
			)
			return
		}
		newCtx := context.WithValue(r.Context(), auth.UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(newCtx))
	})
}
