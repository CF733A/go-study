package auth

import (
    "encoding/json"
    "geo-service/internal/service"
    "geo-service/pkg/responder"
    "net/http"
)

type AuthHandler struct {
    authService service.AuthService
    responder   responder.Responder
}

func NewAuthHandler(authService service.AuthService, responder responder.Responder) *AuthHandler {
    return &AuthHandler{
        authService: authService,
        responder:   responder,
    }
}

// Register godoc
// @Summary Регистрация пользователя
// @Description Создание нового пользователя в системе
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Данные для регистрации"
// @Success 200 {object} responder.Response{data=RegisterResponse} "Успешная регистрация"
// @Failure 400 {object} responder.Response "Ошибка валидации"
// @Failure 500 {object} responder.Response "Внутренняя ошибка сервера"
// @Router /api/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
    var req RegisterRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.responder.BadRequest(w, "invalid json format")
        return
    }
    defer r.Body.Close()

    if err := h.authService.Register(req.Username, req.Password); err != nil {
        if err.Error() == "user already exists" {
            h.responder.Success(w, RegisterResponse{
                Message: "user already exists",
            })
            return
        }
        h.responder.Error(w, http.StatusBadRequest, err.Error())
        return
    }

    h.responder.Success(w, RegisterResponse{
        Message: "user created successfully",
    })
}

// Login godoc
// @Summary Аутентификация пользователя
// @Description Вход в систему и получение JWT токена
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Данные для входа"
// @Success 200 {object} responder.Response{data=LoginResponse} "Успешная аутентификация"
// @Failure 400 {object} responder.Response "Неверные учетные данные"
// @Failure 500 {object} responder.Response "Внутренняя ошибка сервера"
// @Router /api/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    var req LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.responder.BadRequest(w, "invalid json format")
        return
    }
    defer r.Body.Close()

    token, err := h.authService.Login(req.Username, req.Password)
    if err != nil {
        h.responder.Error(w, http.StatusBadRequest, err.Error())
        return
    }

    h.responder.Success(w, LoginResponse{
        Token: token,
    })
}