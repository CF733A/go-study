package home

import (
	"net/http"
	"pet-store/internal/server/http/response"
	"strings"
)

// Home godoc
// @Summary API Home page
// @Description Returns service information in JSON or plain text format based on Accept header
// @Tags info
// @Produce json
// @Produce plain
// @Success 200 {object} HomeResponse "Service information in JSON format"
// @Success 200 {string} string "Service information in plain text format"
// @Router / [get]
// @Example response (application/json)
// {
//   "service": "pet-store",
//   "status": "ok", 
//   "message": "Добро пожаловать в PetStore",
//   "documentation": "http://localhost:8080/swagger/"
// }
// @Example response (text/plain)
// Добро пожаловать на главную страницу PetStore Service
//
// Документация: http://localhost:8080/swagger/
func Home(w http.ResponseWriter, r *http.Request) {
	acceptHeader := r.Header.Get("Accept")

	if strings.Contains(acceptHeader, "application/json") {
		data := NewHomeResponse()
		response.JSONer(w, http.StatusOK, data, nil)
	} else {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Добро пожаловать на главную страницу PetStore Service\n\n" +
			"Документация: http://localhost:8080/swagger/"))
	}
}
