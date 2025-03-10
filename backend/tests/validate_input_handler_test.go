package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"forum/backend/handlers"

	"github.com/stretchr/testify/assert"
)

func TestValidateInputHandler_InvalidMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/validate?username=newUser", nil)
	rr := httptest.NewRecorder()

	handlers.ValidateInputHandler(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
}

func TestValidateInputHandler_InvalidPath(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/invalidpath?username=newUser", nil)
	rr := httptest.NewRecorder()

	handlers.ValidateInputHandler(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestValidateInputHandler_InvalidInput(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/validate", nil)
	rr := httptest.NewRecorder()

	handlers.ValidateInputHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
