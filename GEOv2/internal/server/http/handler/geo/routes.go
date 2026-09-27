package geo

import "github.com/go-chi/chi/v5"

func (h *GeoHandler) RegisterRoutes(r chi.Router) {
    r.Post("/search", h.Search)
    r.Post("/geocode", h.Geocode)
}