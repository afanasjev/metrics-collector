package gauge

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/afanasjev/metrics-collector/internal/handler/testutil"
	"github.com/afanasjev/metrics-collector/internal/repository/memstorage"
)

func TestSetGauge(t *testing.T) {
	testCases := []struct {
		name       string
		value      string
		wantStatus int
		wantStored bool
		wantValue  float64
	}{
		{name: "valid", value: "2.5", wantStatus: http.StatusOK, wantStored: true, wantValue: 2.5},
		{name: "invalid value", value: "NaNnotfloat", wantStatus: http.StatusBadRequest},
	}

	storage := memstorage.GetMemStorage()
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			metricName := testutil.MetricName(t, "handler-gauge")
			req, rr := testutil.RequestWithPathValues(t, http.MethodPost, map[string]string{
				"metricName":  metricName,
				"metricValue": tt.value,
			})

			Set(rr, req)

			if rr.Code != tt.wantStatus {
				t.Fatalf("Set status = %d; want %d", rr.Code, tt.wantStatus)
			}

			if tt.wantStored {
				value, err := storage.GetGauge(metricName)
				if err != nil {
					t.Fatalf("GetGauge failed: %v", err)
				}
				if value != tt.wantValue {
					t.Fatalf("gauge %s = %v, want %v", metricName, value, tt.wantValue)
				}
				return
			}

			if _, err := storage.GetGauge(metricName); err == nil {
				t.Fatalf("gauge %q unexpectedly stored after invalid request", metricName)
			}
		})
	}
}

func TestGetGauge(t *testing.T) {
	t.Run("existing metric", func(t *testing.T) {
		storage := memstorage.GetMemStorage()
		metricName := testutil.MetricName(t, "handler-gauge")
		const expected = 1.234
		if err := storage.SetGauge(metricName, expected); err != nil {
			t.Fatalf("SetGauge failed: %v", err)
		}

		req, rr := testutil.RequestWithPathValues(t, http.MethodGet, map[string]string{
			"metricName": metricName,
		})
		Get(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("Get status = %d; want %d", rr.Code, http.StatusOK)
		}

		if got := strings.TrimSpace(rr.Body.String()); got != fmt.Sprintf("%v", expected) {
			t.Fatalf("Get response = %q; want %q", got, fmt.Sprintf("%v", expected))
		}
	})

	t.Run("missing metric", func(t *testing.T) {
		req, rr := testutil.RequestWithPathValues(t, http.MethodGet, map[string]string{
			"metricName": testutil.MetricName(t, "handler-gauge"),
		})
		Get(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Fatalf("Get status = %d; want %d", rr.Code, http.StatusNotFound)
		}
	})
}
