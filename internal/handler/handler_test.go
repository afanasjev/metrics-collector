package handler

import (
	"testing"
)

func TestGetMetricName(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		expected    string
		expectError bool
	}{
		{
			name:        "valid path with single segment",
			path:        "metricName",
			expected:    "metricName",
			expectError: false,
		},
		{
			name:        "valid path with multiple segments",
			path:        "metricName/value",
			expected:    "metricName",
			expectError: false,
		},
		{
			name:        "path with leading slash",
			path:        "/metricName",
			expected:    "",
			expectError: true,
		},
		{
			name:        "empty path",
			path:        "",
			expected:    "",
			expectError: true,
		},
		{
			name:        "path starting with slash",
			path:        "/test/metric",
			expected:    "",
			expectError: true,
		},
		{
			name:        "path with value",
			path:        "counter/123",
			expected:    "counter",
			expectError: false,
		},
		{
			name:        "path with multiple slashes",
			path:        "gauge/123.45/extra",
			expected:    "gauge",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GetMetricName(tt.path)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for path %q, but got nil", tt.path)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for path %q: %v", tt.path, err)
				}
				if result != tt.expected {
					t.Errorf("Expected metric name %q, got %q", tt.expected, result)
				}
			}
		})
	}
}

