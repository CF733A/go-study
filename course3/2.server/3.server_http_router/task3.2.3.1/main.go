package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Print(".env not found")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := chi.NewRouter()

	r.Get("/", homeHandler)
	r.Get("/1", helloHandler1)
	r.Get("/2", helloHandler2)
	r.Get("/3", helloHandler3)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
	})

	log.Printf("Сервер запущен на порту %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Home page"))
	if err != nil {
		log.Printf("Ошибка записи ответа: %v", err)
	}
}

func helloHandler1(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Hello world"))
	if err != nil {
		log.Printf("Ошибка записи ответа: %v", err)
	}
}

func helloHandler2(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Hello world 2"))
	if err != nil {
		log.Printf("Ошибка записи ответа: %v", err)
	}
}

func helloHandler3(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Hello world 3"))
	if err != nil {
		log.Printf("Ошибка записи ответа: %v", err)
	}
}
