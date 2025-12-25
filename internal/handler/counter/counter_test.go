package counter

import (
	"net/http"
	"strings"
	"testing"

	"github.com/afanasjev/metrics-collector/internal/handler/testutil"
	"github.com/afanasjev/metrics-collector/internal/repository/memstorage"
)

func TestSetCounter(t *testing.T) {
	testCases := []struct {
		name       string
		value      string
		wantStatus int
		wantStored bool
		wantValue  int64
	}{
		{name: "valid", value: "5", wantStatus: http.StatusOK, wantStored: true, wantValue: 5},
		{name: "invalid value", value: "not-an-int", wantStatus: http.StatusBadRequest},
	}

	storage := memstorage.GetMemStorage()
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			metricName := testutil.MetricName(t, "handler-counter")
			req, rr := testutil.RequestWithPathValues(t, http.MethodPost, map[string]string{
				"metricName":  metricName,
				"metricValue": tt.value,
			})

			Set(rr, req)

			if rr.Code != tt.wantStatus {
				t.Fatalf("Set status = %d; want %d", rr.Code, tt.wantStatus)
			}

			if tt.wantStored {
				value, err := storage.GetCounter(metricName)
				if err != nil {
					t.Fatalf("GetCounter failed: %v", err)
				}
				if value != tt.wantValue {
					t.Fatalf("counter %s = %d, want %d", metricName, value, tt.wantValue)
				}
				return
			}

			if _, err := storage.GetCounter(metricName); err == nil {
				t.Fatalf("counter %q unexpectedly stored after invalid request", metricName)
			}
		})
	}
}

func TestGetCounter(t *testing.T) {
	t.Run("existing metric", func(t *testing.T) {
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
			t.Fatalf("Get status = %d; want %d", rr.Code, http.StatusOK)
		}

		if got := strings.TrimSpace(rr.Body.String()); got != "42" {
			t.Fatalf("Get response = %q; want %q", got, "42")
		}
	})

	t.Run("missing metric", func(t *testing.T) {
		req, rr := testutil.RequestWithPathValues(t, http.MethodGet, map[string]string{
			"metricName": testutil.MetricName(t, "handler-counter"),
		})
		Get(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Fatalf("Get status = %d; want %d", rr.Code, http.StatusNotFound)
		}
	})
}
