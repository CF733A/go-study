package domain

import (
	"geo-service/internal/service"
	"geo-service/internal/service/domain/auth"
	"geo-service/internal/service/domain/geo"
	"geo-service/internal/storage"
)

// NewService создает реализацию всех служб
func NewService(repositories *storage.Repository, tokenGen service.TokenGenerator) *service.Service {
    authService := auth.NewAuthService(repositories.User, tokenGen)
    geoService := geo.NewGeoService(repositories.Geocoder)

    return &service.Service{
        Auth: authService,
        Geo:  geoService,
    }
}