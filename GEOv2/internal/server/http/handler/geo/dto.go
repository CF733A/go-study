package geo

// SearchRequestDTO представляет запрос на поиск координат по адресу
type SearchRequestDTO struct {
    Query string `json:"query" example:"Москва, Красная площадь"`
}

// GeocodeRequestDTO представляет запрос на поиск адреса по координатам
type GeocodeRequestDTO struct {
    Lat string `json:"lat" example:"55.7558"`
    Lng string `json:"lng" example:"37.6173"`
}

// SearchResponseDTO представляет ответ для поиска координат
type SearchResponseDTO struct {
    Addresses []AddressResponseDTO `json:"addresses"`
}

type AddressResponseDTO struct {
    Source       string  `json:"source" example:"москва сухонская 11"`
    Result       string  `json:"result" example:"г Москва, ул Сухонская, д 11"`
    PostalCode   string  `json:"postal_code,omitempty" example:"127642"`
    Country      string  `json:"country,omitempty" example:"Россия"`
    Region       string  `json:"region,omitempty" example:"Москва"`
    CityArea     string  `json:"city_area,omitempty" example:"Северо-восточный"`
    CityDistrict string  `json:"city_district,omitempty" example:"Северное Медведково"`
    Street       string  `json:"street,omitempty" example:"Сухонская"`
    House        string  `json:"house,omitempty" example:"11"`
    GeoLat       float64 `json:"geo_lat" example:"55.8782557"`
    GeoLon       float64 `json:"geo_lon" example:"37.65372"`
}

// GeocodeResponseDTO представляет ответ для обратного геокодирования
type GeocodeResponseDTO struct {
    Suggestions []SuggestionDTO `json:"suggestions"`
}

type SuggestionDTO struct {
    Value             string `json:"value" example:"г Москва, ул Сухонская, д 11"`
    UnrestrictedValue string `json:"unrestricted_value" example:"г Москва, ул Сухонская, д 11"`
    Data              SuggestionDataDTO `json:"data,omitempty"`
}

type SuggestionDataDTO struct {
    PostalCode   string  `json:"postal_code,omitempty" example:"127642"`
    Country      string  `json:"country,omitempty" example:"Россия"`
    Region       string  `json:"region,omitempty" example:"Москва"`
    CityArea     string  `json:"city_area,omitempty" example:"Северо-восточный"`
    CityDistrict string  `json:"city_district,omitempty" example:"Северное Медведково"`
    Street       string  `json:"street,omitempty" example:"Сухонская"`
    House        string  `json:"house,omitempty" example:"11"`
    GeoLat       float64 `json:"geo_lat,omitempty" example:"55.8782557"`
    GeoLon       float64 `json:"geo_lon,omitempty" example:"37.65372"`
}