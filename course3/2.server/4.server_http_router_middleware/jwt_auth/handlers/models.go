package handlers

import "github.com/golang-jwt/jwt/v4"

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

// User представляет информацию о пользователе
type User struct {
	Username string `json:"username"` // имя пользователя
	Password string `json:"-"`        // пароль (не показываем в JSON)
}

// AuthRequest - то, что присылает пользователь при регистрации/логине
type AuthRequest struct {
	Username string `json:"username" example:"user123"`
	Password string `json:"password" example:"password123"`
}

// AuthResponse - то, что возвращаем пользователю (JWT токен)
type AuthResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// JWTClaims - дополнительные данные в JWT токене
type JWTClaims struct {
	Username string `json:"username"` // храним имя пользователя в токене
	jwt.RegisteredClaims
}

// LoginResponse - ответ для /api/login (для случаев с ошибками аутентификации)
type LoginResponse struct {
	Success bool   `json:"success" example:"true"`
	Token   string `json:"token,omitempty" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	Error   string `json:"error,omitempty" example:"user not found"`
}