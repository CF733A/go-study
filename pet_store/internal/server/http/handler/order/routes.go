package order

import "github.com/go-chi/chi/v5"

func (h *OrderHandler) RegisterPublicRoutes(r chi.Router) {
	r.Get("/order/{orderId}", h.GetOrderByID)
	r.Delete("/order/{orderId}", h.DeleteOrder)
	r.Post("/order", h.CreateOrder)
}

func (h *OrderHandler) RegisterProtectedRoutes(r chi.Router) {
	r.Get("/store/inventory", h.GetInventory)
}
