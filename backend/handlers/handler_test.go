package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"encoding/json"
	"github.com/Githaiga22/forum/backend/handlers"
	"github.com/Githaiga22/forum/backend/repositories"
	"github.com/Githaiga22/forum/backend/util"
	"github.com/Githaiga22/forum/backend/models"
	"github.com/stretchr/testify/assert"
)

// Mock GetPosts function
func MockGetPosts(db util.Database) ([]models.Post, error) {
	return []models.Post{
		{
			ID:        1,
			UserID:    123,
			Username:  "user1",
			PostTitle: "Post 1",
			Body:      "This is post 1",
			CreatedOn: time.Now().Add(-time.Hour),
			MediaURL:  "http://example.com/media1",
		},
		{
			ID:        2,
			UserID:    124,
			Username:  "user2",
			PostTitle: "Post 2",
			Body:      "This is post 2",
			CreatedOn: time.Now().Add(-2 * time.Hour),
			MediaURL:  "http://example.com/media2",
		},
	}, nil
}

// Mock the SessionStore for testing
var MockSessionStore = map[string]interface{}{
	"mock_session_id": true,
}

// Mock the getSessionID function
func MockGetSessionID(r *http.Request) (string, error) {
	return "mock_session_id", nil
}

// Test HomeHandler
func TestHomeHandler(t *testing.T) {
	// Override the dependencies with mocks
	handlers.GetSessionID = MockGetSessionID
	handlers.SessionStore = MockSessionStore
	repositories.GetPosts = MockGetPosts

	// Test case: Method Not Allowed
	req, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.HomeHandler)
	handler.ServeHTTP(rr, req)

	// Check if method not allowed returns the correct status code
	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("Expected status code %d, got %d", http.StatusMethodNotAllowed, status)
	}

	// Test case: Page not found (wrong URL path)
	req, err = http.NewRequest("GET", "/wrongpath", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Check if page not found returns the correct status code
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, status)
	}

	// Test case: Successful request (valid path and method)
	req, err = http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Check if the status code is OK
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, status)
	}

	// Check if the response contains the mock post details (you can inspect rr.Body)
	var posts []models.Post
	if err := json.NewDecoder(rr.Body).Decode(&posts); err != nil {
		t.Fatal("Failed to decode response:", err)
	}

	// Check if posts were returned and match the mock data
	assert.Len(t, posts, 2)
	assert.Equal(t, "Post 1", posts[0].PostTitle)
	assert.Equal(t, "Post 2", posts[1].PostTitle)
}
