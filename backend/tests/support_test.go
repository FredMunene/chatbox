package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"forum/backend/handlers"
)

func TestCreateSession(t *testing.T) {
	sessionID := handlers.CreateSession()
	if sessionID == "" {
		t.Errorf("Expected non-empty session ID")
	}

	// Ensure session is created in SessionStore
	if _, exists := handlers.SessionStore[sessionID]; !exists {
		t.Errorf("Session ID not found in SessionStore")
	}
}

func TestSetSessionCookie(t *testing.T) {
	rr := httptest.NewRecorder()
	sessionID := handlers.CreateSession()
	handlers.SetSessionCookie(rr, sessionID)

	cookie := rr.Result().Cookies()
	if len(cookie) == 0 {
		t.Errorf("Expected a session cookie to be set")
	}

	if cookie[0].Name != "session_token" || cookie[0].Value != sessionID {
		t.Errorf("Unexpected cookie values: got %v", cookie[0])
	}
}

func TestGetSessionID(t *testing.T) {
	sessionID := handlers.CreateSession()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "session_token", Value: sessionID})

	retrievedSessionID, err := handlers.GetSessionID(req)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if retrievedSessionID != sessionID {
		t.Errorf("Expected session ID %q, got %q", sessionID, retrievedSessionID)
	}
}

func TestSetAndGetSessionData(t *testing.T) {
	sessionID := handlers.CreateSession()
	handlers.SetSessionData(sessionID, "username", "testuser")

	data, err := handlers.GetSessionData(sessionID)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if username, ok := data["username"].(string); !ok || username != "testuser" {
		t.Errorf("Expected 'testuser', got %v", username)
	}
}

func TestIsValidEmail(t *testing.T) {
	validEmails := []string{"test@example.com", "user.name@domain.co", "valid_email123@sub.domain.net"}
	invalidEmails := []string{"invalid-email", "user@com", "@domain.com"}

	for _, email := range validEmails {
		if !handlers.IsValidEmail(email) {
			t.Errorf("Expected valid email for %q", email)
		}
	}

	for _, email := range invalidEmails {
		if handlers.IsValidEmail(email) {
			t.Errorf("Expected invalid email for %q", email)
		}
	}
}
