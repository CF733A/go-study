package auth

import "github.com/go-chi/chi/v5"

func (h *AuthHandler) RegisterRoutes(r chi.Router) {
    r.Post("/register", h.Register)
    r.Post("/login", h.Login)
}