package http

import (
	"pet-store/internal/auth"
	"pet-store/internal/server/http/handler/home"
	"pet-store/internal/server/http/middleware"

	"github.com/go-chi/chi/v5"
	_ "pet-store/internal/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(controllers *Controllers, jwtManager *auth.JWTManager) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.CommonMiddleware...)

	r.Get("/", home.Home)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))


	r.Group(func(r chi.Router) {
		r.Route("/user", func(r chi.Router) {
			controllers.Auth.RegisterPublicRoutes(r)
			controllers.User.RegisterPublicRoutes(r)
		})
		r.Route("/store", func(r chi.Router) {
			controllers.Order.RegisterPublicRoutes(r)
		})

	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(jwtManager))

		controllers.Order.RegisterProtectedRoutes(r)
		controllers.Pet.RegisterProtectedRoutes(r)

	})

	return r
}
