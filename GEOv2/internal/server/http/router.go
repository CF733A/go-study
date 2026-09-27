package http

import (
	"geo-service/internal/config"
	"geo-service/internal/server/http/middleware"
	"geo-service/pkg/logger"
	"geo-service/pkg/responder"

	"github.com/go-chi/chi/v5"
	mid "github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(
	controllers *Controllers,
	resp responder.Responder,
	appLogger *logger.Logger,
	cfg *config.Config,
	authMiddleware *middleware.AuthMiddleware) *chi.Mux {

	r := chi.NewRouter()

	r.Use(mid.Logger)
	r.Use(mid.Recoverer)
	r.Use(mid.RequestID)

	loggerMiddleware := middleware.NewLoggerMiddleware(appLogger)
	r.Use(loggerMiddleware.Middleware)

	r.Get("/", controllers.Home.Home)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	r.Route("/api", func(r chi.Router) {
		controllers.Auth.RegisterRoutes(r)

		r.Route("/address", func(r chi.Router) {
			r.Use(authMiddleware.Middleware)
			controllers.Geo.RegisterRoutes(r)
		})
	})

	return r
}
