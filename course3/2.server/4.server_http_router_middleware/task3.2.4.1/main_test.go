package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRoutes(t *testing.T) {
	r := InitRouter()

	tests := []struct {
		method   string
		path     string
		status   int
		expected string
	}{
		{"GET", "/", 200, "welcome"},
		{"POST", "/", 405, ""},

		{"GET", "/1", 200, "Hello World"},
		{"GET", "/2", 200, "Hello World 1"},
		{"GET", "/notfound", 404, "Not Found"},
		{"POST", "/1", 405, ""},
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
				t.Errorf("For %s %s: expected status %d, got %d",
					test.method, test.path, test.status, rr.Code)
			}

			if test.status == 200 || test.status == 404 {
				if rr.Body.String() != test.expected {
					t.Errorf("For %s %s: expected body '%s', got '%s'",
						test.method, test.path, test.expected, rr.Body.String())
				}
			}
		})
	}
}