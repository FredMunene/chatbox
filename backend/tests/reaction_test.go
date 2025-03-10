package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"forum/backend/handlers"
)

func TestReactionHandler(t *testing.T) {
	reqBody := url.Values{}
	reqBody.Set("reaction", "Like")
	reqBody.Set("post_id", "1")

	req := httptest.NewRequest(http.MethodPost, "/reaction", strings.NewReader(reqBody.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handlers.ReactionHandler(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Expected status %d, got %d", http.StatusSeeOther, rr.Code)
	}
}
