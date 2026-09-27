package auth

// LoginRequest represents user login credentials
// @Description User authentication request payload
type LoginRequest struct {
	Username string `json:"username" example:"john_doe"`  // Username for login
	Password string `json:"password" example:"password123"`  // Password for login
}

// LoginResponse represents successful authentication response
// @Description JWT token response after successful login
type LoginResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`  // JWT access token
	Type  string `json:"type" example:"Bearer"`  // Token type
}