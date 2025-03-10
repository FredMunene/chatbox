package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateSession(t *testing.T) {
	sessionID := CreateSession()
	if sessionID == "" {
		t.Errorf("Expected non-empty session ID")
	}

	// Ensure session is created in SessionStore
	if _, exists := SessionStore[sessionID]; !exists {
		t.Errorf("Session ID not found in SessionStore")
	}
}

func TestSetSessionCookie(t *testing.T) {
	rr := httptest.NewRecorder()
	sessionID := CreateSession()
	SetSessionCookie(rr, sessionID)

	cookie := rr.Result().Cookies()
	if len(cookie) == 0 {
		t.Errorf("Expected a session cookie to be set")
	}

	if cookie[0].Name != "session_token" || cookie[0].Value != sessionID {
		t.Errorf("Unexpected cookie values: got %v", cookie[0])
	}
}

func TestGetSessionID(t *testing.T) {
	sessionID := CreateSession()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "session_token", Value: sessionID})

	retrievedSessionID, err := getSessionID(req)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if retrievedSessionID != sessionID {
		t.Errorf("Expected session ID %q, got %q", sessionID, retrievedSessionID)
	}
}

func TestSetAndGetSessionData(t *testing.T) {
	sessionID := CreateSession()
	SetSessionData(sessionID, "username", "testuser")

	data, err := getSessionData(sessionID)
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
		if !isValidEmail(email) {
			t.Errorf("Expected valid email for %q", email)
		}
	}

	for _, email := range invalidEmails {
		if isValidEmail(email) {
			t.Errorf("Expected invalid email for %q", email)
		}
	}
}
