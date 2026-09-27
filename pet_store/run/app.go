package run

import (
	"fmt"
	"pet-store/config"
	"pet-store/internal/auth"
	"pet-store/internal/migrator"
	Http "pet-store/internal/server/http"
	services "pet-store/internal/services/domain"
	"pet-store/internal/storage/postgres"

	"github.com/go-chi/chi/v5"
)

type App struct {
	Router *chi.Mux
	Config *config.Config
}

func NewApp() (*App, error) {
	cfg := config.LoadConfig()

	db, err := postgres.NewDB(cfg)
	if err != nil {
		return nil, err
	}

	migrator := migrator.NewMigrator()
	if err := migrator.Migrate(cfg); err != nil {
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}

	// repo
	repo := postgres.NewPostgresRepository(db)

	// services
	svc := services.NewServices(repo, cfg)

	// jwt
	jwtManager := auth.NewJWTManager(cfg.JWTSecret)

	// controllers
	controllers := Http.NewControllers(svc)

	// router
	router := Http.NewRouter(controllers, jwtManager)

	// app
	return &App{
		Router: router,
		Config: cfg,
	}, nil
}
