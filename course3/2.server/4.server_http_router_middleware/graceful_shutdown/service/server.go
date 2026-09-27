package service

import (
	"context"
	"log"
	"net/http"
	"time"

	"studentgit.kata.academy/f4ud1991_mail.ru/go-kata/course3/2.server/4.server_http_router_middleware/jwt_auth/config"
)

type MyServer struct {
	server *http.Server
	config *config.Config
}

func InitMyServer(cfg *config.Config, handler http.Handler) *MyServer {
	return &MyServer{
		server: &http.Server{
			Addr:         ":" + cfg.Port,
			Handler:      handler,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
		config: cfg,
	}
}

func (s *MyServer) Serve() error {
	log.Printf("AppServer starting on port %s", s.config.Port)
	return s.server.ListenAndServe()
}

func (s *MyServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
