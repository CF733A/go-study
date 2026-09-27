package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewRouterInit(t *testing.T) {
	r := NewRouterInit()
	if r == nil {
		t.Error("Router init error")
	}
}

func testServerInit(t *testing.T) *httptest.Server {
	r := NewRouterInit()
	server := httptest.NewServer(r)
	t.Cleanup(func() {
		server.Close()
	})
	return server
}

func TestApiHandler(t *testing.T) {
	server := testServerInit(t)

	resp, err := http.Get(server.URL + "/api")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected %d, got %d", http.StatusOK, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != "Hello from API" {
		t.Errorf("expected body 'Hello from API', got %s", string(body))
	}
}

// func TestHugoHandler(t *testing.T){
// 	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.WriteHeader(http.StatusOK)
// 		w.Write([]byte("Test server"))
// 	}))
// 	defer server.Close()


// }