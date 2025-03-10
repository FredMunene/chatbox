package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndexHandler(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		url          string
		expectedCode int
	}{
		{"Invalid Path", "GET", "/wrongpath", http.StatusNotFound},
		{"Invalid Method", "POST", "/home", http.StatusMethodNotAllowed},
		{"Valid Request", "GET", "/home", http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.url, nil)
			rr := httptest.NewRecorder()

			IndexHandler(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)
		})
	}
}
