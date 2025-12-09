package counter

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
	return fmt.Sprintf("handler-counter-%s", strings.ReplaceAll(t.Name(), "/", "-"))
}

func TestSetCounterSuccess(t *testing.T) {
	name := metricNameFromTest(t)
	req := routeRequest(t, http.MethodPost, map[string]string{
		"metricName":  name,
		"metricValue": "5",
	})

	rr := httptest.NewRecorder()
	Set(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("Set returned status %d, want %d", status, http.StatusOK)
	}

	storage := memstorage.GetMemStorage()
	value, err := storage.GetCounter(name)
	if err != nil {
		t.Fatalf("GetCounter failed: %v", err)
	}
	if value != 5 {
		t.Fatalf("counter %s = %d, want 5", name, value)
	}
}

func TestSetCounterBadValue(t *testing.T) {
	name := metricNameFromTest(t)
	req := routeRequest(t, http.MethodPost, map[string]string{
		"metricName":  name,
		"metricValue": "invalid",
	})

	rr := httptest.NewRecorder()
	Set(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Fatalf("Set returned status %d, want %d", status, http.StatusBadRequest)
	}

	if _, err := memstorage.GetMemStorage().GetCounter(name); err == nil {
		t.Fatalf("expected counter %q to be unset after bad request", name)
	}
}

func TestGetCounterOutputsStoredValue(t *testing.T) {
	name := metricNameFromTest(t)
	storage := memstorage.GetMemStorage()
	if err := storage.SetCounter(name, 42); err != nil {
		t.Fatalf("SetCounter failed: %v", err)
	}

	req := routeRequest(t, http.MethodGet, map[string]string{
		"metricName": name,
	})
	rr := httptest.NewRecorder()
	Get(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("Get returned status %d, want %d", status, http.StatusOK)
	}

	if got := strings.TrimSpace(rr.Body.String()); got != "42" {
		t.Fatalf("Get response %q, want %q", got, "42")
	}
}

