package user

import "github.com/go-chi/chi/v5"

func (h *UserHandler) RegisterPublicRoutes(r chi.Router) {
	r.Post("/", h.CreateUser)
	r.Post("/createWithArray", h.CreateUsersWithArray)
	r.Post("/createWithList", h.CreateUsersWithList)
	r.Get("/{username}", h.GetUserByUsername)
	r.Put("/{username}", h.UpdateUser)
	r.Delete("/{username}", h.DeleteUser)
}
