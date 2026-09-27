package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRoute(t *testing.T) {
	r := InitRouter()
	logger = InitLogger()
	if logger == nil{
		t.Fatal("Logger not created")
	}

	tests := []struct {
		method string
		path   string
		status int
		body   string
	}{
		{"GET", "/", 200, "Welcome"},
		{"GET", "/something", 200, "Hello from page 2"},
		{"GET", "/1", 404, "Not Found"},
		{"POST", "/", 405, ""},
	}

	for _, test := range tests {
		t.Run(test.method+test.path, func(*testing.T) {
			req, err := http.NewRequest(test.method, test.path, nil)
			if err != nil {
				t.Fatal(err)
			}

			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			if rec.Code != test.status {
				t.Errorf("For %s %s: expected status %d, got %d",
					test.method, test.path, test.status, rec.Code)
			}

			if test.status == 200 || test.status == 404 {
				if rec.Body.String() != test.body {
					t.Errorf("For %s %s: expected body '%s', got '%s'",
						test.method, test.path, test.body, rec.Body.String())
				}
			}
		})
	}
}
