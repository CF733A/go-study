package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main(){
	err := godotenv.Load()
	if err != nil {
		log.Printf("")
	}

	port := os.Getenv("PORT")
	if port == ""{
		port = "8080"
	}


	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Привет, мир!"))
	})

	log.Printf("Сервер запущен на порту %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}