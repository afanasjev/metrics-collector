package update

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_ServeHTTP(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{
			name:           "POST request returns BadRequest",
			method:         http.MethodPost,
			path:           "",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "GET request returns BadRequest",
			method:         http.MethodGet,
			path:           "",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "PUT request returns BadRequest",
			method:         http.MethodPut,
			path:           "",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "DELETE request returns BadRequest",
			method:         http.MethodDelete,
			path:           "",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "any path returns BadRequest",
			method:         http.MethodPost,
			path:           "some/path",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "http://example.com/", bytes.NewBuffer([]byte{}))
			req.URL.Path = tt.path
			rec := httptest.NewRecorder()

			handler := &Handler{}
			handler.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rec.Code)
			}
		})
	}
}

