package handlers

import (
	"encoding/json"
	"forum/backend/handlers"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// func TestSignupHandler_Success(t *testing.T) {
// 	reqBody := url.Values{}
// 	reqBody.Set("username", "testuser")
// 	reqBody.Set("email", "test@example.com")
// 	reqBody.Set("password", "securepassword")
// 	reqBody.Set("confirmed-password", "securepassword")

// 	req := httptest.NewRequest(http.MethodPost, "/sign-up", strings.NewReader(reqBody.Encode()))
// 	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

// 	rr := httptest.NewRecorder()

// 	handlers.SignupHandler(rr, req)

// 	if rr.Code != http.StatusSeeOther {
// 		t.Errorf("Expected status %d, got %d", http.StatusSeeOther, rr.Code)
// 	}

// 	var response handlers.Response
// 	err := json.NewDecoder(rr.Body).Decode(&response)
// 	if err != nil {
// 		t.Errorf("Failed to decode JSON response: %v", err)
// 	}
// 	if !response.Success {
// 		t.Errorf("Expected success response, got failure")
// 	}
// }

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
