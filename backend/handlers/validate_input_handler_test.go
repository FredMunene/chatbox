package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"forum/backend/util"
	"github.com/stretchr/testify/assert"
)

// Mock database query function
type MockDB struct{}

func (m *MockDB) QueryRow(query string, args ...interface{}) *sql.Row {
	if args[0] == "takenUser" || args[0] == "taken@example.com" {
		return &sql.Row{} // Simulate existing user
	}
	return &sql.Row{} // Simulate no record found
}

func TestValidateInputHandler(t *testing.T) {
	util.Database = &MockDB{} // Inject mock DB

	tests := []struct {
		name       string
		queryParam string
		expected   bool
	}{
		{"Username Available", "username=newUser", true},
		{"Username Taken", "username=takenUser", false},
		{"Email Available", "email=new@example.com", true},
		{"Email Taken", "email=taken@example.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/validate?"+tt.queryParam, nil)
			rr := httptest.NewRecorder()

			ValidateInputHandler(rr, req)

			var response map[string]bool
			err := json.NewDecoder(rr.Body).Decode(&response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, response["available"])
		})
	}
}
