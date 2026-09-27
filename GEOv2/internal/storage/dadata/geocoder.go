package dadata

import (
    "context"
    "errors"
    "geo-service/internal/entities"
    "strconv"

    "github.com/ekomobile/dadata/v2"
    "github.com/ekomobile/dadata/v2/api/clean"
    "github.com/ekomobile/dadata/v2/api/suggest"
    "github.com/ekomobile/dadata/v2/client"
)

type DaDataGeocoder struct {
    cleanApi   *clean.Api
    suggestApi *suggest.Api
}

func NewDaDataGeocoder(apiKey, secretKey string) *DaDataGeocoder {
    creds := &client.Credentials{
        ApiKeyValue:    apiKey,
        SecretKeyValue: secretKey,
    }

    cleanApi := dadata.NewCleanApi(client.WithCredentialProvider(creds))
    suggestApi := dadata.NewSuggestApi(client.WithCredentialProvider(creds))

    return &DaDataGeocoder{
        cleanApi:   cleanApi,
        suggestApi: suggestApi,
    }
}

func (d *DaDataGeocoder) GetCoordinates(req *entities.SearchRequest) (*entities.SearchResponse, error) {
    addresses, err := d.cleanApi.Address(context.Background(), req.Query)
    if err != nil {
        return nil, err
    }

    if len(addresses) == 0 {
        return nil, errors.New("address not found")
    }

    domainAddresses := make([]entities.Address, len(addresses))
    for i, addr := range addresses {
        lat, _ := strconv.ParseFloat(addr.GeoLat, 64)
        lng, _ := strconv.ParseFloat(addr.GeoLon, 64)

        domainAddresses[i] = entities.Address{
            Source:       req.Query,
            Result:       addr.Result,
            PostalCode:   addr.PostalCode,
            Country:      addr.Country,
            Region:       addr.Region,
            CityArea:     addr.CityArea,
            CityDistrict: addr.CityDistrict,
            Street:       addr.Street,
            House:        addr.House,
            GeoLat:       lat,
            GeoLon:       lng,
        }
    }

    return &entities.SearchResponse{
        Addresses: domainAddresses,
    }, nil
}

func (d *DaDataGeocoder) GetAddresses(req *entities.GeocodeRequest) (*entities.GeocodeResponse, error) {
    params := &suggest.GeolocateParams{
        Lat: req.Lat,
        Lon: req.Lng,
    }

    addresses, err := d.suggestApi.GeoLocate(context.Background(), params)
    if err != nil {
        return nil, err
    }

    suggestions := make([]entities.Suggestion, len(addresses))
    for i, addr := range addresses {
        lat, _ := strconv.ParseFloat(addr.Data.GeoLat, 64)
        lon, _ := strconv.ParseFloat(addr.Data.GeoLon, 64)

        suggestions[i] = entities.Suggestion{
            Value:             addr.Value,
            UnrestrictedValue: addr.UnrestrictedValue,
            Data: entities.SuggestionData{
                PostalCode:   addr.Data.PostalCode,
                Country:      addr.Data.Country,
                Region:       addr.Data.Region,
                CityArea:     addr.Data.CityArea,
                CityDistrict: addr.Data.CityDistrict,
                Street:       addr.Data.Street,
                House:        addr.Data.House,
                GeoLat:       lat,
                GeoLon:       lon,
            },
        }
    }

    return &entities.GeocodeResponse{
        Suggestions: suggestions,
    }, nil
}