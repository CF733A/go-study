package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"pet-store/internal/auth"
	"pet-store/internal/server/http/response"
)

type contextKey string

const (
	UserIDKey contextKey = "userID" // - это мне надо для контекста, чтобы передать в заказ
)

func AuthMiddleware(jwtManager *auth.JWTManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				response.JSONer(w, http.StatusUnauthorized, nil, errors.New("authorization header required"))
				return
			}

			parts := strings.Split(authHeader, " ")

			if len(parts) != 2 || parts[0] != "Bearer" {
				response.JSONer(w, http.StatusUnauthorized, nil, errors.New("authorization header format should be: Bearer <token>"))
				return
			}

			tokenString := parts[1]

			if tokenString == "" {
				response.JSONer(w, http.StatusUnauthorized, nil, errors.New("token is empty"))
				return
			}

			userID, err := jwtManager.Validate(tokenString)
			if err != nil {
				response.JSONer(w, http.StatusUnauthorized, nil, errors.New("invalid token: "+err.Error()))
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserIDFromContext(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(UserIDKey).(int64)
	return userID, ok
}