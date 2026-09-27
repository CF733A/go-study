package handlers

import (
	"net/http"
	"strings"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Expected: Bearer <token>", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		claims, err := VerifyJWT(tokenString)
		if err != nil {
			http.Error(w, "Invalid or expired token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		if !UserExists(claims.Username) {
			http.Error(w, "User no longer exists", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
