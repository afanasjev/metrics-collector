package gauge

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/afanasjev/metrics-collector/internal/repository/memstorage"
)

func routeRequest(t *testing.T, method string, params map[string]string) *http.Request {
	t.Helper()
	req := httptest.NewRequest(method, "/", nil)
	for key, value := range params {
		req.SetPathValue(key, value)
	}
	return req
}

func metricNameFromTest(t *testing.T) string {
	t.Helper()
	return fmt.Sprintf("handler-gauge-%s", strings.ReplaceAll(t.Name(), "/", "-"))
}

func TestSetGaugeSuccess(t *testing.T) {
	name := metricNameFromTest(t)
	req := routeRequest(t, http.MethodPost, map[string]string{
		"metricName":  name,
		"metricValue": "2.5",
	})

	rr := httptest.NewRecorder()
	Set(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("Set returned status %d, want %d", status, http.StatusOK)
	}

	storage := memstorage.GetMemStorage()
	value, err := storage.GetGauge(name)
	if err != nil {
		t.Fatalf("GetGauge failed: %v", err)
	}
	if value != 2.5 {
		t.Fatalf("gauge %s = %g, want %g", name, value, 2.5)
	}
}

func TestSetGaugeBadValue(t *testing.T) {
	name := metricNameFromTest(t)
	req := routeRequest(t, http.MethodPost, map[string]string{
		"metricName":  name,
		"metricValue": "NaNnotfloat",
	})

	rr := httptest.NewRecorder()
	Set(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Fatalf("Set returned status %d, want %d", status, http.StatusBadRequest)
	}

	if _, err := memstorage.GetMemStorage().GetGauge(name); err == nil {
		t.Fatalf("expected gauge %q to be unset after bad request", name)
	}
}

func TestGetGaugeOutputsStoredValue(t *testing.T) {
	name := metricNameFromTest(t)
	storage := memstorage.GetMemStorage()
	const expected = 1.234
	if err := storage.SetGauge(name, expected); err != nil {
		t.Fatalf("SetGauge failed: %v", err)
	}

	req := routeRequest(t, http.MethodGet, map[string]string{
		"metricName": name,
	})
	rr := httptest.NewRecorder()
	Get(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("Get returned status %d, want %d", status, http.StatusOK)
	}

	if got := strings.TrimSpace(rr.Body.String()); got != fmt.Sprintf("%v", expected) {
		t.Fatalf("Get response %q, want %q", got, fmt.Sprintf("%v", expected))
	}
}

