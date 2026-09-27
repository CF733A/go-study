package middleware

import (
    "geo-service/pkg/logger"
    "net/http"
    "time"
	"go.uber.org/zap"
)

type LoggerMiddleware struct {
    logger *logger.Logger
}

func NewLoggerMiddleware(logger *logger.Logger) *LoggerMiddleware {
    return &LoggerMiddleware{
        logger: logger,
    }
}

func (m *LoggerMiddleware) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        // Создаем wrapped writer для отслеживания статуса
        wrappedWriter := &responseWriter{
            ResponseWriter: w,
            statusCode:     http.StatusOK,
        }

        next.ServeHTTP(wrappedWriter, r)

        duration := time.Since(start)

        m.logger.Info("HTTP request",
            zap.String("method", r.Method),
            zap.String("path", r.URL.Path),
            zap.Int("status", wrappedWriter.statusCode),
            zap.Duration("duration", duration),
            zap.String("ip", r.RemoteAddr),
        )
    })
}

type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}