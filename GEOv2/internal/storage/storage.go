package storage

import (
    "geo-service/internal/entities"
)

type UserRepository interface {
    Create(user *entities.User) error
    FindByUsername(username string) (*entities.User, error)
}

type Geocoder interface {
    GetCoordinates(req *entities.SearchRequest) (*entities.SearchResponse, error)
    GetAddresses(req *entities.GeocodeRequest) (*entities.GeocodeResponse, error)
}

type Repository struct {
    User    UserRepository
    Geocoder Geocoder
}