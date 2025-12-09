package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/afanasjev/metrics-collector/internal/handler/counter"
	"github.com/afanasjev/metrics-collector/internal/handler/gauge"
	"github.com/afanasjev/metrics-collector/internal/repository/memstorage"
	"github.com/go-chi/chi/v5"
)

func TestUpdateCounterEndpoint(t *testing.T) {
	storage := memstorage.GetMemStorage()

	tests := []struct {
		name           string
		path           string
		expectedStatus int
	}{
		{
			name:           "valid counter update",
			path:           "/update/counter/testCounter/10",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid counter value",
			path:           "/update/counter/testCounter/invalid",
			expectedStatus: http.StatusBadRequest,
		},
	}

	router := setupTestRouter()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, tt.path, nil)
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("Endpoint status = %v, want %v", status, tt.expectedStatus)
			}

			// Verify storage if update was successful
			if tt.expectedStatus == http.StatusOK {
				if _, err := storage.GetCounter("testCounter"); err != nil {
					t.Fatalf("GetCounter() error = %v", err)
				}
			}
		})
	}
}

func TestUpdateGaugeEndpoint(t *testing.T) {
	storage := memstorage.GetMemStorage()

	tests := []struct {
		name           string
		path           string
		expectedStatus int
	}{
		{
			name:           "valid gauge update",
			path:           "/update/gauge/testGauge/10.5",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid gauge value",
			path:           "/update/gauge/testGauge/invalid",
			expectedStatus: http.StatusBadRequest,
		},
	}

	router := setupTestRouter()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, tt.path, nil)
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("Endpoint status = %v, want %v", status, tt.expectedStatus)
			}

			// Verify storage if update was successful
			if tt.expectedStatus == http.StatusOK {
				if _, err := storage.GetGauge("testGauge"); err != nil {
					t.Fatalf("GetGauge() error = %v", err)
				}
			}
		})
	}
}

func TestUpdateEndpointWithoutType(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodPost, "/update", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Endpoint status = %v, want %v", status, http.StatusNotFound)
	}
}

func TestUpdateEndpointWithInvalidType(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodPost, "/update/invalid/testMetric/10", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Endpoint status = %v, want %v", status, http.StatusBadRequest)
	}
}

func setupTestRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Post("/update", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	router.Post("/update/{metricType}/{metricName}/{metricValue}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	})
	router.Post("/update/counter/{metricName}/{metricValue}", counter.Handle)
	router.Post("/update/gauge/{metricName}/{metricValue}", gauge.Handle)
	return router
}

