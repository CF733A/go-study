package main

import (
    _ "geo-service/internal/docs"
    "geo-service/internal/server/http"
    "geo-service/run"
    "log"
)

// @title GEO Service API
// @version 1.0
// @description Сервис для геокодирования адресов с аутентификацией JWT

// @contact.name API Support
// @contact.url http://localhost:8080
// @contact.email support@geo-service.com

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Введите: Bearer {JWT_TOKEN}

func main() {
    app, err := run.NewApp()
    if err != nil {
        log.Fatal("Failed to create app:", err)
    }

    // Создаем и запускаем сервер
    server := http.NewServer(app.Config, app.Router)
    manager := http.NewLifecycleManager(server)
    manager.Run()
}