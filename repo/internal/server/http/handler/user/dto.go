package user

import "time"

// ErrorResponse структура для ошибок API
// @Description Стандартный ответ об ошибке
type ErrorResponse struct {
	Error   string `json:"error" example:"Bad Request"`
	Message string `json:"message" example:"Invalid request parameters"`
	Code    int    `json:"code" example:"400"`
}

// CreateUserRequest для создания пользователя
// @Description Запрос на создание пользователя
type CreateUserRequest struct {
	Name  string `json:"name" example:"Иван Иванов"`
	Email string `json:"email" example:"ivan@example.com"`
}

// UpdateUserRequest для обновления пользователя
// @Description Запрос на обновление пользователя
type UpdateUserRequest struct {
	Name  string `json:"name" example:"Иван Иванов"`
	Email string `json:"email" example:"ivan@example.com"`
}

// UserResponse для ответа API
// @Description Ответ с данными пользователя
type UserResponse struct {
	ID        string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name      string    `json:"name" example:"Иван Иванов"`
	Email     string    `json:"email" example:"ivan@example.com"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserListResponse для списка пользователей
// @Description Ответ со списком пользователей
type UserListResponse struct {
	Users  []UserResponse `json:"users"`
	Total  int            `json:"total" example:"100"`
	Limit  int            `json:"limit" example:"10"`
	Offset int            `json:"offset" example:"0"`
}
