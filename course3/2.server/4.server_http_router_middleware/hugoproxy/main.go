package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := NewRouterInit()
	log.Fatal(http.ListenAndServe(":8080", r))
}

func NewRouterInit() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api", func(r chi.Router){
		r.HandleFunc("/*", apiHandler)
	})
	r.HandleFunc("/*", hugoHandler)

	return r
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello from API"))
}

func hugoHandler(w http.ResponseWriter, r *http.Request) {
	targetURL, err := url.Parse("http://hugo:1313/")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	origDir := proxy.Director
	proxy.Director = func(req *http.Request) {
		origDir(req)
		req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
		req.Host = targetURL.Host
	}

	proxy.ServeHTTP(w, r)
}
