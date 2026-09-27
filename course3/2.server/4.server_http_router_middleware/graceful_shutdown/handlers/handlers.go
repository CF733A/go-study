package handlers

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
)

// HomeHandler godoc
// @Summary Главная страница
// @Description Возвращает приветственное сообщение
// @Tags main
// @Success 200 {string} string "Добро пожаловать на главную страницу"
// @Router / [get]
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Добро пожаловать на главную страницу"))
}

// GeoHandler godoc
// @Summary Геокодирование координат
// @Description Преобразует широту и долготу в адрес
// @Tags geocode
// @Accept json
// @Produce json
// @Param request body GeocodeRequest true "Координаты для геокодирования"
// @Success 200 {array} AddressResponse "Успешный ответ с адресами"
// @Failure 400 {string} string "Неверный формат запроса"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /api/address/geocode [post]
func GeoHandler(w http.ResponseWriter, r *http.Request) {
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

// AdressHandler godoc
// @Summary Поиск адреса
// @Description Ищет адрес по текстовому запросу
// @Tags address
// @Accept json
// @Produce json
// @Param request body SearchRequest true "Поисковый запрос"
// @Success 200 {array} AddressResponse "Успешный ответ с адресами"
// @Failure 400 {string} string "Неверный формат запроса"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /api/address/search [post]
func AdressHandler(w http.ResponseWriter, r *http.Request) {
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

// RegisterHandler godoc
// @Summary Регистрация нового пользователя
// @Description Создает аккаунт пользователя и возвращает JWT токен
// @Tags auth
// @Accept json
// @Produce json
// @Param request body AuthRequest true "Данные для регистрации"
// @Success 200 {object} AuthResponse "Успешная регистрация"
// @Failure 400 {string} string "Неверный запрос"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /api/register [post]
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method", http.StatusMethodNotAllowed)
		return
	}

	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
	}

	if req.Username == "" || req.Password == "" {
		http.Error(w, "Fields 'username', 'password' and 'email' are required", http.StatusBadRequest)
		return
	}

	if UserExists(req.Username){
		http.Error(w, "User already exists", http.StatusBadRequest)
		return
	}

	hashedPass, err := HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	user := User{
		Username: req.Username,
		Password: hashedPass,
	}
	SaveUser(user)

	token, err := GenerateJWT(req.Username)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	response := AuthResponse{Token: token}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}


// LoginHandler godoc
// @Summary Аутентификация пользователя
// @Description Проверяет логин/пароль и возвращает JWT токен
// @Tags auth
// @Accept json
// @Produce json
// @Param request body AuthRequest true "Данные для входа"
// @Success 200 {object} LoginResponse "Успешный вход или ошибка аутентификации"
// @Failure 400 {string} string "Неверный запрос"
// @Failure 405 {string} string "Метод не разрешен"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /api/login [post]
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" {
		http.Error(w, "Fields 'username' and 'password' are required", http.StatusBadRequest)
		return
	}

	user, exists := GetUser(req.Username)
	if !exists {
		response := LoginResponse{
			Success: false,
			Error:   "user not found",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	if !CheckPasswordHash(req.Password, user.Password) {
		response := LoginResponse{
			Success: false,
			Error:   "invalid password",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	token, err := GenerateJWT(req.Username)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	response := LoginResponse{
		Success: true,
		Token:   token,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

