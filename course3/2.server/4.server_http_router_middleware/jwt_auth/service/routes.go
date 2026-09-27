package service

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "studentgit.kata.academy/f4ud1991_mail.ru/go-kata/course3/2.server/4.server_http_router_middleware/jwt_auth/docs"
	"studentgit.kata.academy/f4ud1991_mail.ru/go-kata/course3/2.server/4.server_http_router_middleware/jwt_auth/handlers"
)

func InitRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
	})

	Routes(r)

	return r
}

func Routes(r *chi.Mux) {
	r.Get("/", handlers.HomeHandler)

	r.Post("/api/register", handlers.RegisterHandler)
	r.Post("/api/login", handlers.LoginHandler)

	r.Get("/swagger/*", httpSwagger.WrapHandler)


	r.Route("/api/address", func(r chi.Router) {
		r.Use(handlers.AuthMiddleware)

		r.Route("/search", func(r chi.Router) {
			r.Post("/", handlers.AdressHandler)
		})

		r.Route("/geocode", func(r chi.Router) {
			r.Post("/", handlers.GeoHandler)
		})
	})
}
