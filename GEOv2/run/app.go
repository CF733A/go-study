package run

import (
	"geo-service/internal/config"
	"geo-service/internal/server/http"
	"geo-service/internal/server/http/middleware"
	"geo-service/internal/service/domain"
	"geo-service/internal/storage"
	"geo-service/internal/storage/dadata"
	"geo-service/internal/storage/memory"
	"geo-service/pkg/jwt"
	"geo-service/pkg/logger"
	"geo-service/pkg/responder"

	"github.com/go-chi/chi/v5"
)

type App struct {
    Router *chi.Mux
    Config *config.Config
}

func NewApp() (*App, error) {
    // config init
    cfg := config.LoadConfig()
    
    // utils init
    jsonResponder := responder.NewJSONResponder()
    appLogger := logger.New()
    tokenizer := jwt.NewJWTTokenizer(cfg.JWTsecret)

    // repo init
    userRepo := memory.NewUserRepository()
    geocoder := dadata.NewDaDataGeocoder(cfg.DadataApiKey, cfg.DadataSecret)
    
    repositories := &storage.Repository{
        User:    userRepo,
        Geocoder: geocoder,
    }
    
    // service init
    services := domain.NewService(repositories, tokenizer)

    // controllers init
    controllers := http.NewControllers(services, jsonResponder)

    // middleware init
    authMiddleware := middleware.NewAuthMiddleware(services.Auth, jsonResponder)
    
    // router init
    router := http.NewRouter(controllers, jsonResponder, appLogger, cfg, authMiddleware)
    
    return &App{
        Router: router,
        Config: cfg,
    }, nil
}