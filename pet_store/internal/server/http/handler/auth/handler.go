package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"pet-store/internal/server/http/response"
	services "pet-store/internal/services"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login godoc
// @Summary User login
// @Description Authenticate user with credentials and return JWT token
// @Tags authentication
// @Accept json
// @Produce json
// @Param request body LoginRequest true "User login credentials"
// @Success 200 {object} LoginResponse "Successfully authenticated"
// @Failure 400 {object} response.Response "Invalid input data"
// @Failure 401 {object} response.Response "Invalid username or password"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /user/login [post]
// @Example request
// {
//   "username": "john_doe",
//   "password": "password123"
// }
// @Example response
// {
//   "status": true,
//   "data": {
//     "token": "example.jwt.token",
//     "type": "Bearer"
//   }
// }
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.JSONer(w, http.StatusMethodNotAllowed, nil, errors.New("method not allowed"))
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSONer(w, http.StatusBadRequest, nil, errors.New("invalid JSON format"))
		return
	}

	if req.Username == "" || req.Password == "" {
		response.JSONer(w, http.StatusBadRequest, nil, errors.New("username and password are required"))
		return
	}

	token, err := h.authService.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		response.JSONer(w, http.StatusUnauthorized, nil, err)
		return
	}

	responseData := LoginResponse{
		Token: token,
		Type:  "Bearer",
	}

	response.JSONer(w, http.StatusOK, responseData, nil)
}

// Logout godoc
// @Summary User logout
// @Description Revoke user's JWT token to end session
// @Tags authentication
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response "Successfully logged out"
// @Failure 400 {object} response.Response "Invalid or missing token"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /user/logout [get]
// @Example response
// {
//   "status": true,
//   "message": "successfully logged out"
// }
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.JSONer(w, http.StatusMethodNotAllowed, nil, errors.New("method not allowed"))
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		response.JSONer(w, http.StatusBadRequest, nil, errors.New("authorization header required"))
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		response.JSONer(w, http.StatusBadRequest, nil, errors.New("authorization header format should be: Bearer <token>"))
		return
	}

	token := parts[1]
	if token == "" {
		response.JSONer(w, http.StatusBadRequest, nil, errors.New("token is empty"))
		return
	}

	if err := h.authService.Logout(r.Context(), token); err != nil {
		response.JSONer(w, http.StatusBadRequest, nil, err)
		return
	}

	response.JSONer(w, http.StatusOK, map[string]string{"message": "successfully logged out"}, nil)
}