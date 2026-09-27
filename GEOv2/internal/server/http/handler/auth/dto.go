package auth

// RegisterRequest представляет запрос на регистрацию
type RegisterRequest struct {
    Username string `json:"username" example:"john_doe"`
    Password string `json:"password" example:"securepassword123"`
}

// LoginRequest представляет запрос на аутентификацию
type LoginRequest struct {
    Username string `json:"username" example:"john_doe"`
    Password string `json:"password" example:"securepassword123"`
}

// RegisterResponse представляет ответ при успешной регистрации
type RegisterResponse struct {
    Message string `json:"message" example:"user created successfully"`
}

// LoginResponse представляет ответ при успешной аутентификации
type LoginResponse struct {
    Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}