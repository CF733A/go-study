package auth

import "github.com/go-chi/chi/v5"

func (h *AuthHandler) RegisterPublicRoutes(r chi.Router) {
    r.Route("/", func(r chi.Router) {
        r.Post("/login", h.Login)
        r.Get("/logout", h.Logout)
    })
}