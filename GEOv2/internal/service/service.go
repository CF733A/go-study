package service

import (
    "geo-service/internal/entities"
)

// TokenGenerator интерфейс для работы с JWT токенами
type TokenGenerator interface {
    Generate(user *entities.User) (string, error)
    Validate(token string) (*entities.User, error)
}

// AuthService интерфейс для аутентификации
type AuthService interface {
    Register(username, password string) error
    Login(username, password string) (string, error)
    ValidateToken(token string) (*entities.User, error)
}

// GeoService интерфейс для геокодирования
type GeoService interface {
    AddressToCoordinates(req *entities.SearchRequest) (*entities.SearchResponse, error)
    CoordinatesToAddress(req *entities.GeocodeRequest) (*entities.GeocodeResponse, error)
}

// Service объединяет все службы
type Service struct {
    Auth AuthService
    Geo  GeoService
}