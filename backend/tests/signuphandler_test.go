package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"forum/backend/handlers"
)

func TestSignupHandler_MismatchedPasswords(t *testing.T) {
	reqBody := url.Values{}
	reqBody.Set("username", "testuser")
	reqBody.Set("email", "test@example.com")
	reqBody.Set("password", "securepassword")
	reqBody.Set("confirmed-password", "differentpassword")

	req := httptest.NewRequest(http.MethodPost, "/sign-up", strings.NewReader(reqBody.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handlers.SignupHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var response handlers.Response
	err := json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Errorf("Failed to decode JSON response: %v", err)
	}
	if response.Success {
		t.Errorf("Expected failure response due to mismatched passwords")
	}
}

func TestSignupHandler_InvalidFormFields(t *testing.T) {
	reqBody := url.Values{}
	reqBody.Set("username", "")
	reqBody.Set("email", "invalid-email")
	reqBody.Set("password", "short")
	reqBody.Set("confirmed-password", "short")

	req := httptest.NewRequest(http.MethodPost, "/sign-up", strings.NewReader(reqBody.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handlers.SignupHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var response handlers.Response
	err := json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Errorf("Failed to decode JSON response: %v", err)
	}
	if response.Success {
		t.Errorf("Expected failure response due to invalid form fields")
	}
}

func TestSignupHandler_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodPut, "/sign-up", nil)
	rr := httptest.NewRecorder()

	handlers.SignupHandler(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, rr.Code)
	}
}

func TestSignupHandler_PageNotFound(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/invalid-path", nil)
	rr := httptest.NewRecorder()

	handlers.SignupHandler(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, rr.Code)
	}
}


