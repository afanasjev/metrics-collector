package counter

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/afanasjev/metrics-collector/internal/repository/memstorage"
)

func TestHandle(t *testing.T) {
	tests := []struct {
		name           string
		metricName     string
		metricValue    string
		expectedStatus int
	}{
		{
			name:           "valid counter update",
			metricName:     "testCounter",
			metricValue:    "10",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid metric value",
			metricName:     "testCounter",
			metricValue:    "invalid",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "negative counter value",
			metricName:     "testCounter",
			metricValue:    "-5",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "zero counter value",
			metricName:     "testCounter",
			metricValue:    "0",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "large counter value",
			metricName:     "testCounter",
			metricValue:    "9223372036854775807",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset storage for each test
			storage := memstorage.GetMemStorage()

			req := httptest.NewRequest(http.MethodPost, "/update/counter/"+tt.metricName+"/"+tt.metricValue, nil)
			req.SetPathValue("metricName", tt.metricName)
			req.SetPathValue("metricValue", tt.metricValue)

			rr := httptest.NewRecorder()
			Handle(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("Handle() status = %v, want %v", status, tt.expectedStatus)
			}

			// If status is OK, verify the counter was updated
			if tt.expectedStatus == http.StatusOK {
				value, exists := storage.GetCounter(tt.metricName)
				if !exists {
					t.Error("Counter should exist after successful update")
				}
				if value == 0 && tt.metricValue != "0" {
					t.Errorf("Counter value should not be zero for non-zero input")
				}
			}
		})
	}
}

func TestHandleMultipleUpdates(t *testing.T) {
	storage := memstorage.GetMemStorage()
	metricName := "multiUpdateCounter"

	// First update
	req1 := httptest.NewRequest(http.MethodPost, "/update/counter/"+metricName+"/10", nil)
	req1.SetPathValue("metricName", metricName)
	req1.SetPathValue("metricValue", "10")
	rr1 := httptest.NewRecorder()
	Handle(rr1, req1)

	if rr1.Code != http.StatusOK {
		t.Errorf("First update failed with status %v", rr1.Code)
	}

	value1, _ := storage.GetCounter(metricName)
	if value1 != 10 {
		t.Errorf("After first update, counter = %v, want 10", value1)
	}

	// Second update (should increment)
	req2 := httptest.NewRequest(http.MethodPost, "/update/counter/"+metricName+"/5", nil)
	req2.SetPathValue("metricName", metricName)
	req2.SetPathValue("metricValue", "5")
	rr2 := httptest.NewRecorder()
	Handle(rr2, req2)

	if rr2.Code != http.StatusOK {
		t.Errorf("Second update failed with status %v", rr2.Code)
	}

	value2, _ := storage.GetCounter(metricName)
	if value2 != 15 {
		t.Errorf("After second update, counter = %v, want 15", value2)
	}
}

