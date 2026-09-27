// @title Geo Service API
// @version 1.0
// @description API для работы с адресами и геокодированием
// @host localhost:8080
// @BasePath /
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"studentgit.kata.academy/f4ud1991_mail.ru/go-kata/course3/2.server/4.server_http_router_middleware/jwt_auth/config"
	"studentgit.kata.academy/f4ud1991_mail.ru/go-kata/course3/2.server/4.server_http_router_middleware/jwt_auth/handlers"
	"studentgit.kata.academy/f4ud1991_mail.ru/go-kata/course3/2.server/4.server_http_router_middleware/jwt_auth/service"
)

func main() {
	cfg := config.LoadConfig()

	handlers.InitAuth(cfg)

	r := service.InitRouter()

	server := service.InitMyServer(cfg, r)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		if err := server.Serve(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-stop
	log.Println("Received shutdown signal")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	} else {
		log.Println("Server stopped gracefully")
	}
}
