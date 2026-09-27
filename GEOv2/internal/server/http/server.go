package http

import (
    "context"
    "geo-service/internal/config"
    "log"
    "net/http"
    "time"
	"github.com/go-chi/chi/v5"
)

type Server struct {
    server *http.Server
    config *config.Config
}

func NewServer(cfg *config.Config, router *chi.Mux) *Server {
    return &Server{
        server: &http.Server{
            Addr:         ":" + cfg.Port,
            Handler:      router,
            ReadTimeout:  10 * time.Second,
            WriteTimeout: 10 * time.Second,
        },
        config: cfg,
    }
}

func (s *Server) Start() error {
    log.Printf("Server starting on port %s", s.config.Port)
    return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
    return s.server.Shutdown(ctx)
}