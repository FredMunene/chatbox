package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"forum/backend/handlers"
	"forum/backend/models"
	"forum/backend/repositories"
)

func TestGetAllPostsAPI(t *testing.T) {
	// Mock database
	var mockDB *sql.DB

	// Create a request
	r, err := http.NewRequest("GET", "/api/posts", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder
	w := httptest.NewRecorder()

	// Call the handler
	handler := handlers.GetAllPostsAPI(mockDB)
	handler.ServeHTTP(w, r)

	// Check status code
	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Check response content type
	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", w.Header().Get("Content-Type"))
	}

	// Check if response is valid JSON
	var posts []models.Post
	if err := json.NewDecoder(w.Body).Decode(&posts); err != nil {
		t.Errorf("response is not valid JSON: %v", err)
	}
}

func TestFilterPosts(t *testing.T) {
	// Create a request with a query parameter
	r, err := http.NewRequest("GET", "/filter?filter=created", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder
	w := httptest.NewRecorder()

	// Call the handler
	http.HandlerFunc(handlers.FilterPosts).ServeHTTP(w, r)

	// Check status code
	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Check if body contains expected text
	if !strings.Contains(w.Body.String(), "An Unexpected Error Occurred") {
		t.Errorf("expected error message in response, got %s", w.Body.String())
	}
}
