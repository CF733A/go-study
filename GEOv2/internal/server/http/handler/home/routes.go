package home

import "github.com/go-chi/chi/v5"

func (h *HomeHandler) RegisterRoutes(r chi.Router) {
    r.Get("/", h.Home)
}