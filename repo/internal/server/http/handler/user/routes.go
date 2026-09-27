package user

import "github.com/go-chi/chi/v5"

func (u *UserHandler) RegisterRoutes(r chi.Router) {
	r.Post("/users", u.CreateUser)
	r.Get("/users", u.ListUsers)
	r.Get("/users/{id}", u.GetUser)
	r.Put("/users/{id}", u.UpdateUser)
	r.Delete("/users/{id}", u.DeleteUser)
}
