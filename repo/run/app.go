package run

import (
	"repo/internal/config"
	"repo/internal/server/http"
	"repo/internal/service"
	user_serv "repo/internal/service/domain/user"
	"repo/internal/storage"
	"repo/internal/storage/postgres"
	"repo/internal/storage/postgres/user"

	"github.com/go-chi/chi/v5"
)

type App struct {
	Router *chi.Mux
	Config *config.Config
	DB     *storage.Repository
}

func NewApp() (*App, error) {
	// config
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	// db
	db, err := postgres.NewDB(cfg)
	if err != nil {
		return nil, err
	}

	// postgres
	userRepo := user.NewUserRepo(db)

	// storage
	storage := storage.NewRepository(userRepo)

	// services
	usrService := user_serv.NewUserService(storage.User)
	service := service.NewService(usrService)

	// controllers
	controllers := http.NewControllers(service)

	// router
	router := http.NewRouter(*controllers)

	return &App{
		Router: router,
		Config: cfg,
		DB:     storage,
	}, nil
}
