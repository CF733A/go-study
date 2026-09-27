package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	dadata "github.com/ekomobile/dadata/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

    httpSwagger "github.com/swaggo/http-swagger"

	_ "studentgit.kata.academy/f4ud1991_mail.ru/go-kata/course3/2.server/4.server_http_router_middleware/geoservice_swagger/docs"
	
)

// SearchRequest представляет тело запроса для поиска адреса.
type SearchRequest struct {
	Query string `json:"query" example:"москва сухонская 11"` // Строка поиска адреса
}

// GeocodeRequest представляет тело запроса для геокодирования.
type GeocodeRequest struct {
	Lat string `json:"lat" example:"55.7558"` // Широта
	Lng string `json:"lng" example:"37.6173"` // Долгота
}

// AddressResponse представляет успешный ответ от API.
type AddressResponse struct {
	Source       string `json:"source" example:"dadata"`                                     // Источник данных
	Result       string `json:"result" example:"г Москва, ул Сухонская, д 11"`               // Полный адрес
	PostalCode   string `json:"postal_code" example:"127642"`                                // Почтовый индекс
	Country      string `json:"country" example:"Россия"`                                    // Страна
	Region       string `json:"region" example:"Москва"`                                     // Регион
	CityArea     string `json:"city_area" example:"Северо-восточный административный округ"` // Район города
	CityDistrict string `json:"city_district" example:"район Северное Медведково"`           // Внутригородской район
	Street       string `json:"street" example:"улица Сухонская"`                            // Улица
	House        string `json:"house" example:"11"`                                          // Дом
	GeoLat       string `json:"geo_lat" example:"55.877857"`                                 // Географическая широта
	GeoLon       string `json:"geo_lon" example:"37.653672"`                                 // Географическая долгота
}

// @title Геосервис API
// @version 1.0
// @description Микросервис для поиска и геокодирования адресов через DaData.
func main() {
	r := RouterInit()

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal()
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))

}

func RouterInit() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/address/search", func(r chi.Router) {
		r.Post("/", newAdressHandle)
	})

	r.Route("/api/address/geocode", func(r chi.Router) {
		r.Post("/", newGeoHandle)
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Добро пожаловать на главную страницу!\nАдреса для пост запроса /api/address/search и /api/address/geocode"))
	})

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), // URL для локальной разработки
	))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
	})

	return r
}

// newAdressHandle обрабатывает поиск адреса по строке.
// @Summary Поиск адреса
// @Description Поиск адресов по строке запроса с использованием DaData.
// @Tags addresses
// @Accept  json
// @Produce  json
// @Param   request body SearchRequest true "Поисковый запрос"
// @Success 200 {array} AddressResponse "Успешный ответ"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /api/address/search [post]
func newAdressHandle(w http.ResponseWriter, r *http.Request) {
	apiKey := os.Getenv("DADATA_API_KEY")
	secretKey := os.Getenv("DADATA_SECRET_KEY")

	if apiKey == "" || secretKey == "" {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Dadata credentials not configured. Please set DADATA_API_KEY and DADATA_SECRET_KEY environment variables",
		})
		return
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Only POST method is allowed"))
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte("Content-Type must be application/json"))
		return
	}

	var request SearchRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid JSON format"))
		return
	}

	if request.Query == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Field 'query' is required"))
		return
	}

	api := dadata.NewCleanApi()

	result, err := api.Address(context.Background(), request.Query)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to get address from DaData: " + err.Error(),
		})
		return
	}

	var response []AddressResponse
	for _, addr := range result {
		response = append(response, AddressResponse{
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
		})
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// newGeoHandle обрабатывает обратное геокодирование (координаты в адрес).
// @Summary Обратное геокодирование
// @Description Получение адреса по географическим координатам (широте и долготе) с использованием DaData.
// @Tags geocode
// @Accept  json
// @Produce  json
// @Param   request body GeocodeRequest true "Географические координаты"
// @Success 200 {array} AddressResponse "Успешный ответ"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /api/address/geocode [post]
func newGeoHandle(w http.ResponseWriter, r *http.Request) {
	url := "https://suggestions.dadata.ru/suggestions/api/4_1/rs/geolocate/address"

	apikey := os.Getenv("DADATA_API_KEY")
	if apikey == "" {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("API key not configured"))
		return
	}

	var request GeocodeRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid JSON"))
		return
	}

	if request.Lat == "" || request.Lng == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Fields 'lat' and 'lng' are required"))
		return
	}

	lat, err := strconv.ParseFloat(request.Lat, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid latitude format"))
		return
	}

	lng, err := strconv.ParseFloat(request.Lng, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid longitude format"))
		return
	}

	type Coordinates struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	}

	coords := Coordinates{
		Lat: lat,
		Lon: lng,
	}

	jsonData, err := json.Marshal(coords)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error creating request body"))
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to create request"))
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Token "+apikey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("API request failed"))
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error reading response"))
		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("DaData API error: %s", string(body))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("DaData API error"))
		return
	}

	var daDataResponse struct {
		Suggestions []struct {
			Value             string `json:"value"`
			UnrestrictedValue string `json:"unrestricted_value"`
			Data              struct {
				Source       string `json:"source"`
				Result       string `json:"result"`
				PostalCode   string `json:"postal_code"`
				Country      string `json:"country"`
				Region       string `json:"region"`
				CityArea     string `json:"city_area"`
				CityDistrict string `json:"city_district"`
				Street       string `json:"street"`
				House        string `json:"house"`
				GeoLat       string `json:"geo_lat"`
				GeoLon       string `json:"geo_lon"`
			} `json:"data"`
		} `json:"suggestions"`
	}

	if err := json.Unmarshal(body, &daDataResponse); err != nil {
		log.Printf("Failed to parse response: %s, body: %s", err.Error(), string(body))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to parse response"))
		return
	}

	var response []AddressResponse
	for _, suggestion := range daDataResponse.Suggestions {
		response = append(response, AddressResponse{
			Source:       suggestion.Data.Source,
			Result:       suggestion.Data.Result,
			PostalCode:   suggestion.Data.PostalCode,
			Country:      suggestion.Data.Country,
			Region:       suggestion.Data.Region,
			CityArea:     suggestion.Data.CityArea,
			CityDistrict: suggestion.Data.CityDistrict,
			Street:       suggestion.Data.Street,
			House:        suggestion.Data.House,
			GeoLat:       suggestion.Data.GeoLat,
			GeoLon:       suggestion.Data.GeoLon,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
