package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Githaiga22/forum/backend/handlers"
)

func TestLogoutHandler(t *testing.T) {
	req, err := http.NewRequest("POST", "/logout", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rr := httptest.NewRecorder()
	http.HandlerFunc(handlers.LogoutHandler).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusSeeOther)
	}
}
