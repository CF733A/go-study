package geo

import (
    "geo-service/internal/entities"
    "geo-service/internal/storage"
)

type GeoServiceImpl struct {
    geocoder storage.Geocoder
}

func NewGeoService(geocoder storage.Geocoder) *GeoServiceImpl {
    return &GeoServiceImpl{
        geocoder: geocoder,
    }
}

func (g *GeoServiceImpl) AddressToCoordinates(req *entities.SearchRequest) (*entities.SearchResponse, error) {
    return g.geocoder.GetCoordinates(req)
}

func (g *GeoServiceImpl) CoordinatesToAddress(req *entities.GeocodeRequest) (*entities.GeocodeResponse, error) {
    return g.geocoder.GetAddresses(req)
}