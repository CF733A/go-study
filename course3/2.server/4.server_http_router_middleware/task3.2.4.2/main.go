package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

var logger *zap.Logger

func main() {
	logger = InitLogger()
	defer logger.Sync()

	r := InitRouter()

	logger.Info("server started on port :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func InitLogger() (*zap.Logger){
	var err error
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("Failed to initialize logger: ", err)
	}

	return logger
}

func InitRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Use(Log)
	
	r.Get("/", HandleHome)
	r.Get("/something", HandleSomething)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		logger.Warn("404",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
	)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
	})

	return r
}

func Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		logger.Info("Входящий запрос",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("ip", r.RemoteAddr),
		)
		next.ServeHTTP(w, r)

		logger.Info("Запрос обработан",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("ip", r.RemoteAddr),
			zap.Duration("duration", time.Since(start)),

		)
	})
}

func HandleHome(w http.ResponseWriter, r *http.Request){
	w.WriteHeader(http.StatusOK)
	logger.Info("Обработка главной страницы")
	w.Write([]byte("Welcome"))
}

func HandleSomething(w http.ResponseWriter, r *http.Request){
	w.WriteHeader(http.StatusOK)
	logger.Info("Обработка второстепенной страницы")
	w.Write([]byte("Hello from page 2"))
}