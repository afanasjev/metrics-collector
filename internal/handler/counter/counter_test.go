package counter

import (
	"net/http"
	"strings"
	"testing"

	"github.com/afanasjev/metrics-collector/internal/handler/testutil"
	"github.com/afanasjev/metrics-collector/internal/repository/memstorage"
)

func TestSet(t *testing.T) {
	t.Run("valid request writes value", func(t *testing.T) {
		metricName := testutil.MetricName(t, "handler-counter")
		req, rr := testutil.RequestWithPathValues(t, http.MethodPost, map[string]string{
			"metricName":  metricName,
			"metricValue": "5",
		})

		Set(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
		}

		value, err := memstorage.GetMemStorage().GetCounter(metricName)
		if err != nil {
			t.Fatalf("GetCounter failed: %v", err)
		}
		if value != 5 {
			t.Fatalf("counter %s = %d, want %d", metricName, value, 5)
		}
	})

	t.Run("accumulates values", func(t *testing.T) {
		metricName := testutil.MetricName(t, "handler-counter")
		storage := memstorage.GetMemStorage()
		if err := storage.SetCounter(metricName, 3); err != nil {
			t.Fatalf("preparation SetCounter failed: %v", err)
		}

		req, rr := testutil.RequestWithPathValues(t, http.MethodPost, map[string]string{
			"metricName":  metricName,
			"metricValue": "2",
		})
		Set(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
		}

		value, err := storage.GetCounter(metricName)
		if err != nil {
			t.Fatalf("GetCounter failed: %v", err)
		}
		if value != 5 {
			t.Fatalf("counter %s = %d, want %d", metricName, value, 5)
		}
	})

	t.Run("invalid value returns bad request", func(t *testing.T) {
		metricName := testutil.MetricName(t, "handler-counter")
		req, rr := testutil.RequestWithPathValues(t, http.MethodPost, map[string]string{
			"metricName":  metricName,
			"metricValue": "abc",
		})

		Set(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}

		if _, err := memstorage.GetMemStorage().GetCounter(metricName); err == nil {
			t.Fatalf("expected missing counter %q", metricName)
		}
	})
}

func TestGet(t *testing.T) {
	t.Run("existing metric returns value", func(t *testing.T) {
		storage := memstorage.GetMemStorage()
		metricName := testutil.MetricName(t, "handler-counter")
		if err := storage.SetCounter(metricName, 42); err != nil {
			t.Fatalf("SetCounter failed: %v", err)
		}

		req, rr := testutil.RequestWithPathValues(t, http.MethodGet, map[string]string{
			"metricName": metricName,
		})
		Get(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
		}

		if got := strings.TrimSpace(rr.Body.String()); got != "42" {
			t.Fatalf("body = %q, want %q", got, "42")
		}
	})

	t.Run("missing metric returns 404", func(t *testing.T) {
		req, rr := testutil.RequestWithPathValues(t, http.MethodGet, map[string]string{
			"metricName": testutil.MetricName(t, "handler-counter"),
		})
		Get(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusNotFound)
		}
	})
}
