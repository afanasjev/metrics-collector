package counter

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/afanasjev/metrics-collector/internal/handler/testutil"
	"github.com/afanasjev/metrics-collector/internal/repository/memstorage"
)

func TestSet(t *testing.T) {
	storage := memstorage.GetMemStorage()

	t.Run("valid request accumulates values", func(t *testing.T) {
		t.Parallel()
		metricName := testutil.MetricName(t, "handler-counter-set")
		for _, value := range []string{"3", "4"} {
			req, rr := testutil.RequestWithPathValues(t, http.MethodPost, map[string]string{
				"metricName":  metricName,
				"metricValue": value,
			})

			Set(rr, req)

			if rr.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
			}
		}

		got, err := storage.GetCounter(metricName)
		if err != nil {
			t.Fatalf("GetCounter failed: %v", err)
		}
		if got != 7 {
			t.Fatalf("counter %s = %d, want %d", metricName, got, 7)
		}
	})

	t.Run("invalid value returns bad request", func(t *testing.T) {
		t.Parallel()
		metricName := testutil.MetricName(t, "handler-counter-set-invalid")
		req, rr := testutil.RequestWithPathValues(t, http.MethodPost, map[string]string{
			"metricName":  metricName,
			"metricValue": "not-an-integer",
		})
		Set(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}

		if _, err := storage.GetCounter(metricName); err == nil {
			t.Fatalf("expected missing counter %q", metricName)
		}
	})

	t.Run("invalid value does not change existing counter", func(t *testing.T) {
		t.Parallel()
		metricName := testutil.MetricName(t, "handler-counter-set-invalid-keeps")
		if err := storage.SetCounter(metricName, 5); err != nil {
			t.Fatalf("preparation SetCounter failed: %v", err)
		}

		req, rr := testutil.RequestWithPathValues(t, http.MethodPost, map[string]string{
			"metricName":  metricName,
			"metricValue": "oops",
		})
		Set(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}

		got, err := storage.GetCounter(metricName)
		if err != nil {
			t.Fatalf("GetCounter failed: %v", err)
		}
		if got != 5 {
			t.Fatalf("counter %s = %d, want %d", metricName, got, 5)
		}
	})
}

func TestGet(t *testing.T) {
	storage := memstorage.GetMemStorage()

	t.Run("existing metric returns value", func(t *testing.T) {
		t.Parallel()
		metricName := testutil.MetricName(t, "handler-counter-get")
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
		t.Parallel()
		req, rr := testutil.RequestWithPathValues(t, http.MethodGet, map[string]string{
			"metricName": testutil.MetricName(t, "handler-counter-get"),
		})
		Get(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusNotFound)
		}
	})

	t.Run("write failure returns 500", func(t *testing.T) {
		t.Parallel()
		metricName := testutil.MetricName(t, "handler-counter-get")
		if err := storage.SetCounter(metricName, 1); err != nil {
			t.Fatalf("SetCounter failed: %v", err)
		}

		req, _ := testutil.RequestWithPathValues(t, http.MethodGet, map[string]string{
			"metricName": metricName,
		})
		writer := testutil.NewFailingResponseWriter(errors.New("write failure"))

		Get(writer, req)

		if writer.Status() != http.StatusInternalServerError {
			t.Fatalf("status = %d, want %d", writer.Status(), http.StatusInternalServerError)
		}
	})
}
