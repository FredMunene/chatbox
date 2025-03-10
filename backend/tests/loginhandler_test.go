package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"forum/backend/handlers"
)

func TestLoginHandler(t *testing.T) {
	reqBody := []byte(`{"email": "test@example.com", "password": "password123"}`)
	req, err := http.NewRequest("POST", "/sign-in", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	http.HandlerFunc(handlers.LoginHandler).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var res handlers.Response
	if err := json.Unmarshal(rr.Body.Bytes(), &res); err != nil {
		t.Errorf("Could not parse response body: %v", err)
	}

	if !res.Success {
		t.Errorf("Expected success but got false")
	}
}
