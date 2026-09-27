package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
)

func TestRoutes(t *testing.T) {
	// Создаем роутер такой же, как в main()
	r := chi.NewRouter()
	r.Get("/", homeHandler)
	r.Get("/1", helloHandler1)
	r.Get("/2", helloHandler2)
	r.Get("/3", helloHandler3)
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
	})

	// Тестируем все маршруты
	tests := []struct {
		method   string
		path     string
		status   int
		expected string
	}{
		{"GET", "/", 200, "Home page"},
		{"GET", "/1", 200, "Hello world"},
		{"GET", "/2", 200, "Hello world 2"},
		{"GET", "/3", 200, "Hello world 3"},
		{"GET", "/nonexistent", 404, "Not Found"},
		{"POST", "/", 405, ""},
	}

	for _, test := range tests {
		t.Run(test.method+test.path, func(t *testing.T) {
			req, err := http.NewRequest(test.method, test.path, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			if rr.Code != test.status {
				t.Errorf("Expected status %d, got %d", test.status, rr.Code)
			}

			if rr.Body.String() != test.expected {
				t.Errorf("Expected body '%s', got '%s'", test.expected, rr.Body.String())
			}
		})
	}
}

func TestHomeHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(homeHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "Home page"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
