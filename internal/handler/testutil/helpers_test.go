package testutil

import (
	"errors"
	"net/http"
	"strings"
	"testing"
)

func TestMetricNameProducesUniqueAndScopedValue(t *testing.T) {
	first := MetricName(t, "metric")
	second := MetricName(t, "metric")

	if first == second {
		t.Fatalf("MetricName should produce unique values, got %q", first)
	}
	if !strings.Contains(first, "metric-"+strings.ReplaceAll(t.Name(), "/", "-")) {
		t.Fatalf("MetricName(%q) should include test name, got %q", t.Name(), first)
	}
}

func TestRequestWithPathValues(t *testing.T) {
	req, rr := RequestWithPathValues(t, http.MethodPost, map[string]string{
		"metricName":  "cpu",
		"metricValue": "10",
	})

	if req.Method != http.MethodPost {
		t.Fatalf("method = %q, want %q", req.Method, http.MethodPost)
	}
	if got := req.PathValue("metricName"); got != "cpu" {
		t.Fatalf("metricName = %q, want %q", got, "cpu")
	}
	if got := req.PathValue("metricValue"); got != "10" {
		t.Fatalf("metricValue = %q, want %q", got, "10")
	}
	if rr == nil {
		t.Fatalf("expected non-nil response recorder")
	}
}

func TestNewFailingResponseWriter(t *testing.T) {
	expectedErr := errors.New("write failed")
	w := NewFailingResponseWriter(expectedErr)

	w.WriteHeader(http.StatusAccepted)
	if got := w.Status(); got != http.StatusAccepted {
		t.Fatalf("Status = %d, want %d", got, http.StatusAccepted)
	}

	if w.Header() == nil {
		t.Fatalf("Header should not be nil")
	}
	_, err := w.Write([]byte("payload"))
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Write error = %v, want %v", err, expectedErr)
	}
}
