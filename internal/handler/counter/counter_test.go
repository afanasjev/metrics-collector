package counter

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetCounterValue(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		expected    int64
		expectError bool
	}{
		{
			name:        "valid path with value",
			path:        "counterName/123",
			expected:    123,
			expectError: false,
		},
		{
			name:        "valid path with negative value should return error",
			path:        "counterName/-456",
			expected:    0,
			expectError: true,
		},
		{
			name:        "valid path with zero",
			path:        "counterName/0",
			expected:    0,
			expectError: false,
		},
		{
			name:        "path without value",
			path:        "counterName",
			expected:    0,
			expectError: true,
		},
		{
			name:        "empty path",
			path:        "",
			expected:    0,
			expectError: true,
		},
		{
			name:        "path with invalid value",
			path:        "counterName/abc",
			expected:    0,
			expectError: true,
		},
		{
			name:        "path with float value",
			path:        "counterName/123.45",
			expected:    0,
			expectError: true,
		},
		{
			name:        "path with large value",
			path:        "counterName/9223372036854775807",
			expected:    9223372036854775807,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := getCounterValue(tt.path)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for path %q, but got nil", tt.path)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for path %q: %v", tt.path, err)
				}
				if result != tt.expected {
					t.Errorf("Expected value %d, got %d", tt.expected, result)
				}
			}
		})
	}
}

func TestHandler_ServeHTTP(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{
			name:           "valid POST request",
			method:         http.MethodPost,
			path:           "testCounter/123",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "GET request should return NotImplemented",
			method:         http.MethodGet,
			path:           "testCounter/123",
			expectedStatus: http.StatusNotImplemented,
		},
		{
			name:           "PUT request should return NotImplemented",
			method:         http.MethodPut,
			path:           "testCounter/123",
			expectedStatus: http.StatusNotImplemented,
		},
		{
			name:           "DELETE request should return NotImplemented",
			method:         http.MethodDelete,
			path:           "testCounter/123",
			expectedStatus: http.StatusNotImplemented,
		},
		{
			name:           "POST request without metric name",
			method:         http.MethodPost,
			path:           "",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "POST request without value",
			method:         http.MethodPost,
			path:           "testCounter",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "POST request with invalid value",
			method:         http.MethodPost,
			path:           "testCounter/abc",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "POST request with negative value should return BadRequest",
			method:         http.MethodPost,
			path:           "testCounter/-10",
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

func TestHandler_ServeHTTP_Integration(t *testing.T) {
	// Тест на реальное сохранение значения
	req := httptest.NewRequest(http.MethodPost, "http://example.com/", bytes.NewBuffer([]byte{}))
	req.URL.Path = "testCounter/100"
	rec := httptest.NewRecorder()

	handler := &Handler{}
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
	}

	// Повторный запрос с тем же именем должен увеличить значение
	req2 := httptest.NewRequest(http.MethodPost, "http://example.com/", bytes.NewBuffer([]byte{}))
	req2.URL.Path = "testCounter/50"
	rec2 := httptest.NewRecorder()

	handler.ServeHTTP(rec2, req2)

	if rec2.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rec2.Code)
	}
}

