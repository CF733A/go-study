package middleware

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

var CommonMiddleware = []func(http.Handler) http.Handler{
	middleware.RequestID,      // Добавляет ID к каждому запросу
	middleware.RealIP,         // Определение реального IP
	middleware.Logger,         // Логирование запросов
	middleware.Recoverer,      // Восстановление после panic
	middleware.CleanPath,      // Очистка пути от лишних слэшей
}