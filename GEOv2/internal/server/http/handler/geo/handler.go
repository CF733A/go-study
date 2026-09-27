package geo

import (
    "encoding/json"
    "geo-service/internal/entities"
    "geo-service/internal/service"
    "geo-service/pkg/responder"
    "net/http"
)

type GeoHandler struct {
    geoService service.GeoService
    responder  responder.Responder
}

func NewGeoHandler(geoService service.GeoService, responder responder.Responder) *GeoHandler {
    return &GeoHandler{
        geoService: geoService,
        responder:  responder,
    }
}

// Search godoc
// @Summary Поиск координат по адресу
// @Description Прямое геокодирование: преобразование адреса в координаты
// @Tags geo
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body SearchRequestDTO true "Адрес для поиска"
// @Success 200 {object} responder.Response{data=SearchResponseDTO} "Успешный поиск"
// @Failure 400 {object} responder.Response "Ошибка валидации"
// @Failure 401 {object} responder.Response "Неавторизован"
// @Failure 500 {object} responder.Response "Ошибка геокодирования"
// @Router /api/address/search [post]
func (h *GeoHandler) Search(w http.ResponseWriter, r *http.Request) {
    var req SearchRequestDTO
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.responder.BadRequest(w, "invalid json format")
        return
    }
    defer r.Body.Close()

    if req.Query == "" {
        h.responder.BadRequest(w, "empty query")
        return
    }

    domainReq := &entities.SearchRequest{
        Query: req.Query,
    }

    result, err := h.geoService.AddressToCoordinates(domainReq)
    if err != nil {
        h.responder.InternalError(w, "geocoding error: "+err.Error())
        return
    }

    responseDTO := h.toSearchResponseDTO(result)
    h.responder.Success(w, responseDTO)
}

// Geocode godoc
// @Summary Поиск адреса по координатам
// @Description Обратное геокодирование: преобразование координат в адрес
// @Tags geo
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body GeocodeRequestDTO true "Координаты для поиска"
// @Success 200 {object} responder.Response{data=GeocodeResponseDTO} "Успешное геокодирование"
// @Failure 400 {object} responder.Response "Ошибка валидации"
// @Failure 401 {object} responder.Response "Неавторизован"
// @Failure 500 {object} responder.Response "Ошибка геокодирования"
// @Router /api/address/geocode [post]
func (h *GeoHandler) Geocode(w http.ResponseWriter, r *http.Request) {
    var req GeocodeRequestDTO
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.responder.BadRequest(w, "invalid json format")
        return
    }
    defer r.Body.Close()

    if req.Lat == "" || req.Lng == "" {
        h.responder.BadRequest(w, "lat and lng cannot be empty")
        return
    }

    domainReq := &entities.GeocodeRequest{
        Lat: req.Lat,
        Lng: req.Lng,
    }

    result, err := h.geoService.CoordinatesToAddress(domainReq)
    if err != nil {
        h.responder.InternalError(w, "reverse geocoding failed: "+err.Error())
        return
    }

    responseDTO := h.toGeocodeResponseDTO(result)
    h.responder.Success(w, responseDTO)
}

// Преобразование Entity в DTO для поиска
func (h *GeoHandler) toSearchResponseDTO(entity *entities.SearchResponse) SearchResponseDTO {
    addresses := make([]AddressResponseDTO, len(entity.Addresses))
    for i, addr := range entity.Addresses {
        addresses[i] = AddressResponseDTO{
            Source:       addr.Source,
            Result:       addr.Result,
            PostalCode:   addr.PostalCode,
            Country:      addr.Country,
            Region:       addr.Region,
            CityArea:     addr.CityArea,
            CityDistrict: addr.CityDistrict,
            Street:       addr.Street,
            House:        addr.House,
            GeoLat:       addr.GeoLat,
            GeoLon:       addr.GeoLon,
        }
    }
    return SearchResponseDTO{Addresses: addresses}
}

// Преобразование Entity в DTO для геокодирования
func (h *GeoHandler) toGeocodeResponseDTO(entity *entities.GeocodeResponse) GeocodeResponseDTO {
    suggestions := make([]SuggestionDTO, len(entity.Suggestions))
    for i, suggestion := range entity.Suggestions {
        suggestions[i] = SuggestionDTO{
            Value:             suggestion.Value,
            UnrestrictedValue: suggestion.UnrestrictedValue,
            Data: SuggestionDataDTO{
                PostalCode:   suggestion.Data.PostalCode,
                Country:      suggestion.Data.Country,
                Region:       suggestion.Data.Region,
                CityArea:     suggestion.Data.CityArea,
                CityDistrict: suggestion.Data.CityDistrict,
                Street:       suggestion.Data.Street,
                House:        suggestion.Data.House,
                GeoLat:       suggestion.Data.GeoLat,
                GeoLon:       suggestion.Data.GeoLon,
            },
        }
    }
    return GeocodeResponseDTO{Suggestions: suggestions}
}