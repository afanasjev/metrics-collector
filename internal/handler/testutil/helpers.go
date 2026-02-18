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

// NewFailingResponseWriter produces a ResponseWriter that always returns an error when written to.
func NewFailingResponseWriter(writeErr error) *failingResponseWriter {
	return &failingResponseWriter{
		header:   make(http.Header),
		writeErr: writeErr,
	}
}

type failingResponseWriter struct {
	header     http.Header
	statusCode int
	writeErr   error
}

func (f *failingResponseWriter) Header() http.Header {
	return f.header
}

func (f *failingResponseWriter) Write([]byte) (int, error) {
	return 0, f.writeErr
}

func (f *failingResponseWriter) WriteHeader(statusCode int) {
	f.statusCode = statusCode
}

func (f *failingResponseWriter) Status() int {
	return f.statusCode
}
