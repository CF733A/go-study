package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := InitRouter()
	log.Fatal(http.ListenAndServe(":8080", r))
}

func InitRouter()*chi.Mux{
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", handleRoute1)
	r.Get("/1", handleRoute2)
	r.Get("/2", handleRoute3)
	
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
	})

	return r
}

func handleRoute1(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("welcome"))
}

func handleRoute2(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello World"))
}

func handleRoute3(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello World 1"))
}