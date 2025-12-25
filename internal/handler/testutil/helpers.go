package testutil

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
)

var metricCounter uint64

// MetricName returns a unique metric name scoped to the current test.
func MetricName(t *testing.T, prefix string) string {
	t.Helper()
	id := atomic.AddUint64(&metricCounter, 1)
	safeTestName := strings.ReplaceAll(t.Name(), "/", "-")
	return fmt.Sprintf("%s-%s-%d", prefix, safeTestName, id)
}

// RequestWithPathValues creates an HTTP request with path values set on it
// using the net/http.Request helpers that Go 1.22+ exposes.
func RequestWithPathValues(t *testing.T, method string, values map[string]string) (*http.Request, *httptest.ResponseRecorder) {
	t.Helper()
	req := httptest.NewRequest(method, "/", nil)
	for key, value := range values {
		req.SetPathValue(key, value)
	}
	return req, httptest.NewRecorder()
}
