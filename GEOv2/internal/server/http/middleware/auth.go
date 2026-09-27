package middleware

import (
	"geo-service/internal/service"
	"geo-service/pkg/responder"
	"net/http"
	"strings"
)

type AuthMiddleware struct {
	authService service.AuthService
	responder   responder.Responder
}

func NewAuthMiddleware(authService service.AuthService, resp responder.Responder) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
		responder:   resp,
	}
}

func (a *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			a.responder.Unauthorized(w, "authorization header required")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			a.responder.Unauthorized(w, "authorization header format should be: Bearer <token>")
			return
		}

		tokenString := parts[1]

		_, err := a.authService.ValidateToken(tokenString)
		if err != nil {
			a.responder.Unauthorized(w, "invalid token: "+err.Error())
			return
		}

		next.ServeHTTP(w, r)
	})
}