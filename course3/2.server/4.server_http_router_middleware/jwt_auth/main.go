// @title Geo Service API
// @version 1.0
// @description API для работы с адресами и геокодированием
// @host localhost:8080
// @BasePath /
package main

import (
	"log"
	"net/http"

	"studentgit.kata.academy/f4ud1991_mail.ru/go-kata/course3/2.server/4.server_http_router_middleware/jwt_auth/config"
	"studentgit.kata.academy/f4ud1991_mail.ru/go-kata/course3/2.server/4.server_http_router_middleware/jwt_auth/handlers"
	"studentgit.kata.academy/f4ud1991_mail.ru/go-kata/course3/2.server/4.server_http_router_middleware/jwt_auth/service"
)

func main() {
	cfg := config.LoadConfig()

	handlers.InitAuth(cfg)

	r := service.InitRouter()

	port := cfg.Port
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
