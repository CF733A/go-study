package main

import (
	"log"
	_ "repo/internal/docs"
	"repo/internal/migrator"
	"repo/internal/server/http"
	"repo/run"
)

// @title User Management API
// @version 1.0
// @description REST API для управления пользователями с PostgreSQL и автоматическими миграциями

// @contact.name API Support
// @contact.url http://localhost:8080
// @contact.email support@example.com

// @host localhost:8080
// @BasePath /
// @schemes http

func main() {
	app, err := run.NewApp()
	if err != nil {
		log.Fatal("Failed to create app:", err)
	}

	migrator := migrator.NewMigrator()
	if err := migrator.Migrate(*app.Config); err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}
	log.Println("Database migrations applied successfully")

	server := http.NewServer(app.Config.ServerPort, app.Router)
	manager := http.NewLifecycleManager(server)
	manager.Run()
}
