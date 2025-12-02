package gauge

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetGaugeValue(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		expected    float64
		expectError bool
	}{
		{
			name:        "valid path with integer value",
			path:        "gaugeName/123",
			expected:    123.0,
			expectError: false,
		},
		{
			name:        "valid path with float value",
			path:        "gaugeName/123.45",
			expected:    123.45,
			expectError: false,
		},
		{
			name:        "valid path with negative value",
			path:        "gaugeName/-456.78",
			expected:    -456.78,
			expectError: false,
		},
		{
			name:        "valid path with zero",
			path:        "gaugeName/0",
			expected:    0.0,
			expectError: false,
		},
		{
			name:        "valid path with zero float",
			path:        "gaugeName/0.0",
			expected:    0.0,
			expectError: false,
		},
		{
			name:        "path without value",
			path:        "gaugeName",
			expected:    0.0,
			expectError: true,
		},
		{
			name:        "empty path",
			path:        "",
			expected:    0.0,
			expectError: true,
		},
		{
			name:        "path with invalid value",
			path:        "gaugeName/abc",
			expected:    0.0,
			expectError: true,
		},
		{
			name:        "path with scientific notation",
			path:        "gaugeName/1.23e-4",
			expected:    0.000123,
			expectError: false,
		},
		{
			name:        "path with very small value",
			path:        "gaugeName/0.000001",
			expected:    0.000001,
			expectError: false,
		},
		{
			name:        "path with very large value",
			path:        "gaugeName/999999999.999",
			expected:    999999999.999,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := getGaugeValue(tt.path)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for path %q, but got nil", tt.path)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for path %q: %v", tt.path, err)
				}
				// Используем небольшую погрешность для сравнения float
				epsilon := 0.0001
				diff := result - tt.expected
				if diff < 0 {
					diff = -diff
				}
				if diff > epsilon {
					t.Errorf("Expected value %f, got %f (diff: %f)", tt.expected, result, diff)
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
			name:           "valid POST request with integer",
			method:         http.MethodPost,
			path:           "testGauge/123",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "valid POST request with float",
			method:         http.MethodPost,
			path:           "testGauge/123.45",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "GET request should return NotImplemented",
			method:         http.MethodGet,
			path:           "testGauge/123",
			expectedStatus: http.StatusNotImplemented,
		},
		{
			name:           "PUT request should return NotImplemented",
			method:         http.MethodPut,
			path:           "testGauge/123",
			expectedStatus: http.StatusNotImplemented,
		},
		{
			name:           "DELETE request should return NotImplemented",
			method:         http.MethodDelete,
			path:           "testGauge/123",
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
			path:           "testGauge",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "POST request with invalid value",
			method:         http.MethodPost,
			path:           "testGauge/abc",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "POST request with negative value",
			method:         http.MethodPost,
			path:           "testGauge/-10.5",
			expectedStatus: http.StatusOK,
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
	req.URL.Path = "testGauge/100.5"
	rec := httptest.NewRecorder()

	handler := &Handler{}
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
	}

	// Повторный запрос должен перезаписать значение
	req2 := httptest.NewRequest(http.MethodPost, "http://example.com/", bytes.NewBuffer([]byte{}))
	req2.URL.Path = "testGauge/200.75"
	rec2 := httptest.NewRecorder()

	handler.ServeHTTP(rec2, req2)

	if rec2.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rec2.Code)
	}
}

