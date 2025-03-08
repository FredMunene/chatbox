package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/Githaiga22/forum/backend/handlers"
	"github.com/Githaiga22/forum/backend/repositories"
	"github.com/Githaiga22/forum/backend/util"
	"github.com/Githaiga22/forum/backend/models"
	"github.com/stretchr/testify/assert"
	"fmt"
)


func MockGetUserByEmail(email string) (models.User, error) {
	return models.User{ID: 1, Email: "user@example.com", Name: "Test User"}, nil
}


func MockGetSessionData(cookie string) (map[string]interface{}, error) {
	return map[string]interface{}{"userEmail": "user@example.com"}, nil
}

// Helper function to make HTTP request and check response
func makeRequest(t *testing.T, method, url string, handler http.HandlerFunc) *httptest.ResponseRecorder {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	return rr
}

// Test IndexHandler
func TestIndexHandler(t *testing.T) {
	// Setup Mocks
	handlers.GetSessionID = MockGetSessionID
	repositories.GetUserByEmail = MockGetUserByEmail
	repositories.GetPosts = MockGetPosts

	// Test cases setup
	tests := []struct {
		name               string
		method             string
		url                string
		expectedCode       int
		mockSessionID      func(*http.Request) (string, error)
		mockSessionData    func() (map[string]interface{}, error)
	}{
		{
			name:         "Invalid Path",
			method:       "GET",
			url:          "/wrongpath",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "Invalid Method",
			method:       "POST",
			url:          "/home",
			expectedCode: http.StatusMethodNotAllowed,
		},
		{
			name:         "Invalid Session ID",
			method:       "GET",
			url:          "/home",
			expectedCode: http.StatusSeeOther,
			mockSessionID: func(r *http.Request) (string, error) {
				return "", fmt.Errorf("invalid session")
			},
		},
		{
			name:         "Valid Request",
			method:       "GET",
			url:          "/home",
			expectedCode: http.StatusOK,
			mockSessionID: func(r *http.Request) (string, error) {
				return "mock_session_id", nil
			},
			mockSessionData: func() (map[string]interface{}, error) {
				return map[string]interface{}{"userEmail": "user@example.com"}, nil
			},
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Override mocks
			if tt.mockSessionID != nil {
				handlers.GetSessionID = tt.mockSessionID
			} else {
				handlers.GetSessionID = MockGetSessionID
			}

			if tt.mockSessionData != nil {
				handlers.GetSessionData = tt.mockSessionData
			} else {
				handlers.GetSessionData = MockGetSessionData
			}

			// Perform the request
			rr := makeRequest(t, tt.method, tt.url, handlers.IndexHandler)

			// Assert status code
			assert.Equal(t, tt.expectedCode, rr.Code)

			// If valid request, verify posts are returned
			if tt.expectedCode == http.StatusOK {
				var posts []models.Post
				if err := json.NewDecoder(rr.Body).Decode(&posts); err != nil {
					t.Fatal("Failed to decode response:", err)
				}
				assert.Len(t, posts, 1)
				assert.Equal(t, "Post 1", posts[0].PostTitle)
			}
		})
	}
}
