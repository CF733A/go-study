package home

import (
	"geo-service/pkg/responder"
	"net/http"
	"strings"
)

type HomeHandler struct {
	responder responder.Responder
}

func NewHomeHandler(responder responder.Responder) *HomeHandler {
	return &HomeHandler{
		responder: responder,
	}
}

// Home godoc
// @Summary Главная страница API
// @Description Возвращает информацию о сервисе в JSON или текстовом формате
// @Tags info
// @Produce json
// @Produce plain
// @Success 200 {object} HomeResponse "Информация о сервисе в JSON"
// @Success 200 {string} string "Информация о сервисе в текстовом формате"
// @Router / [get]
func (h *HomeHandler) Home(w http.ResponseWriter, r *http.Request) {
	acceptHeader := r.Header.Get("Accept")

	if strings.Contains(acceptHeader, "application/json") {
		data := NewHomeResponse()
		h.responder.Success(w, data)
	} else {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Добро пожаловать на главную страницу GEO Service\n\n" +
			"Доступные эндпойнты:\n" +
			"- POST /api/register - регистрация\n" +
			"- POST /api/login - авторизация\n" +
			"- POST /api/address/search - поиск координат (требует аутентификацию)\n" +
			"- POST /api/address/geocode - обратное геокодирование (требует аутентификацию)\n\n" +
			"Документация: http://localhost:8080/swagger/"))
	}
}
