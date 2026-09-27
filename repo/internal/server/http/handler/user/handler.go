package user

import (
	"encoding/json"
	"log"
	"net/http"
	"repo/internal/service"
	"repo/internal/storage"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(us service.UserService) *UserHandler {
	return &UserHandler{
		userService: us,
	}
}

// CreateUser создает нового пользователя
// @Summary Создать пользователя
// @Description Создает нового пользователя с указанными именем и email
// @Tags users
// @Accept json
// @Produce json
// @Param user body CreateUserRequest true "Данные пользователя"
// @Success 201 {object} user.User "Пользователь создан"
// @Failure 400 {object} ErrorResponse "Неверные данные"
// @Failure 409 {object} ErrorResponse "Пользователь с таким email уже существует"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/users [post]
func (u *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "wrong json format", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Name) == "" {
		http.Error(w, "empty name", http.StatusBadRequest)
		return
	}

	if len(req.Name) < 2 {
		http.Error(w, "min 2 chars", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Email) == "" {
		http.Error(w, "empty email", http.StatusBadRequest)
		return
	}

	user, err := u.userService.CreateUser(r.Context(), req.Name, req.Email)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			sendJSON(w, http.StatusOK, ErrorResponse{
				Error:   "error",
				Message: "user already exists",
				Code:    http.StatusConflict,
			})
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError) // 500
		}
		return
	}

	sendJSON(w, http.StatusCreated, user)
}

// GetUser возвращает пользователя по ID
// @Summary Получить пользователя
// @Description Возвращает пользователя по его UUID
// @Tags users
// @Produce json
// @Param id path string true "UUID пользователя"
// @Success 200 {object} user.User "Данные пользователя"
// @Failure 400 {object} ErrorResponse "Неверный UUID"
// @Failure 404 {object} ErrorResponse "Пользователь не найден"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/users/{id} [get]
func (u *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if _, err := uuid.Parse(id); err != nil {
		http.Error(w, "wrong id", http.StatusBadRequest)
		return
	}

	user, err := u.userService.GetUser(r.Context(), id)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	sendJSON(w, http.StatusOK, user)
}

// UpdateUser обновляет данные пользователя
// @Summary Обновить пользователя
// @Description Обновляет имя и email пользователя
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "UUID пользователя"
// @Param user body UpdateUserRequest true "Новые данные пользователя"
// @Success 200 {object} user.User "Обновленные данные пользователя"
// @Failure 400 {object} ErrorResponse "Неверные данные"
// @Failure 404 {object} ErrorResponse "Пользователь не найден"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/users/{id} [put]
func (u *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if _, err := uuid.Parse(id); err != nil {
		http.Error(w, "wrong id", http.StatusBadRequest)
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "wrong json format", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Name) == "" {
		http.Error(w, "empty name", http.StatusBadRequest)
		return
	}

	user, err := u.userService.UpdateUser(r.Context(), id, req.Name, req.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendJSON(w, http.StatusOK, user)
}

// DeleteUser удаляет пользователя
// @Summary Удалить пользователя
// @Description Выполняет мягкое удаление пользователя (soft delete)
// @Tags users
// @Produce json
// @Param id path string true "UUID пользователя"
// @Success 200 {object} map[string]string "Сообщение об успешном удалении"
// @Failure 400 {object} ErrorResponse "Неверный UUID"
// @Failure 404 {object} ErrorResponse "Пользователь не найден"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/users/{id} [delete]
func (u *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if _, err := uuid.Parse(id); err != nil {
		http.Error(w, "wrong id", http.StatusBadRequest)
		return
	}

	if err := u.userService.DeleteUser(r.Context(), id); err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "user successfully deleted",
	})
}

// ListUsers возвращает список пользователей
// @Summary Получить список пользователей
// @Description Возвращает список пользователей с пагинацией
// @Tags users
// @Produce json
// @Param limit query int false "Лимит (макс. 100)" default(10) minimum(1) maximum(100)
// @Param offset query int false "Смещение" default(0) minimum(0)
// @Success 200 {object} UserListResponse "Список пользователей"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/users [get]
func (u *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	if limit <= 0 || limit > 100 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	cond := storage.Conditions{
		Limit:  limit,
		Offset: offset,
	}

	users, total, err := u.userService.ListUsers(r.Context(), cond)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userList := UserListResponse{
		Users:  make([]UserResponse, len(users)),
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}

	for i, usr := range users {
		userList.Users[i] = UserResponse{
			ID:        usr.ID,
			Name:      usr.Name,
			Email:     usr.Email,
			CreatedAt: usr.CreatedAt,
			UpdatedAt: usr.UpdatedAt,
		}
	}

	sendJSON(w, http.StatusOK, userList)
}

func sendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		log.Printf("JSON encode error: %v", err)
		return
	}
}
