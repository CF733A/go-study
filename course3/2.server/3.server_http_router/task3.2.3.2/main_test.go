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
		{"GET", "/", 200, "Home Page"},
		{"POST", "/", 405, ""},

		{"GET", "/group1/", 200, "Group 1"},
		{"GET", "/group1/1", 200, "Group 1 Привет, мир 1"},
		{"GET", "/group1/2", 200, "Group 1 Привет, мир 2"},
		{"GET", "/group1/3", 200, "Group 1 Привет, мир 3"},
		{"POST", "/group1/", 405, ""},

		{"GET", "/group2/", 200, "Group 2"},
		{"GET", "/group2/1", 200, "Group 2 Привет, мир 1"},
		{"GET", "/group2/2", 200, "Group 2 Привет, мир 2"},
		{"GET", "/group2/3", 200, "Group 2 Привет, мир 3"},
		{"POST", "/group2/", 405, ""},

		{"GET", "/group3/", 200, "Group 3"},
		{"GET", "/group3/1", 200, "Group 3 Привет, мир 1"},
		{"GET", "/group3/2", 200, "Group 3 Привет, мир 2"},
		{"GET", "/group3/3", 200, "Group 3 Привет, мир 3"},
		{"POST", "/group3/", 405, ""},

		{"GET", "/nonexistent", 404, "Not Found"},
		{"POST", "/nonexistent", 404, "Not Found"},
		{"GET", "/group1/nonexistent", 404, "Not Found"},
		{"GET", "/invalid/path", 404, "Not Found"},
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