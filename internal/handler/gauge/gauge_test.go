package gauge

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
			name:           "valid gauge update",
			metricName:     "testGauge",
			metricValue:    "10.5",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid metric value",
			metricName:     "testGauge",
			metricValue:    "invalid",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "negative gauge value",
			metricName:     "testGauge",
			metricValue:    "-5.3",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "zero gauge value",
			metricName:     "testGauge",
			metricValue:    "0.0",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "large gauge value",
			metricName:     "testGauge",
			metricValue:    "123456789.987654321",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "scientific notation",
			metricName:     "testGauge",
			metricValue:    "1.5e10",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := memstorage.GetMemStorage()

			req := httptest.NewRequest(http.MethodPost, "/update/gauge/"+tt.metricName+"/"+tt.metricValue, nil)
			req.SetPathValue("metricName", tt.metricName)
			req.SetPathValue("metricValue", tt.metricValue)

			rr := httptest.NewRecorder()
			Handle(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("Handle() status = %v, want %v", status, tt.expectedStatus)
			}

			// If status is OK, verify the gauge was updated
			if tt.expectedStatus == http.StatusOK {
				if _, err := storage.GetGauge(tt.metricName); err != nil {
					t.Fatalf("GetGauge() error = %v", err)
				}
			}
		})
	}
}

func TestHandleMultipleUpdates(t *testing.T) {
	storage := memstorage.GetMemStorage()
	metricName := "multiUpdateGauge"

	// First update
	req1 := httptest.NewRequest(http.MethodPost, "/update/gauge/"+metricName+"/10.5", nil)
	req1.SetPathValue("metricName", metricName)
	req1.SetPathValue("metricValue", "10.5")
	rr1 := httptest.NewRecorder()
	Handle(rr1, req1)

	if rr1.Code != http.StatusOK {
		t.Errorf("First update failed with status %v", rr1.Code)
	}

	value1, err := storage.GetGauge(metricName)
	if err != nil {
		t.Fatalf("GetGauge() after first update = %v", err)
	}
	if value1 != 10.5 {
		t.Errorf("After first update, gauge = %v, want 10.5", value1)
	}

	// Second update (should overwrite)
	req2 := httptest.NewRequest(http.MethodPost, "/update/gauge/"+metricName+"/20.3", nil)
	req2.SetPathValue("metricName", metricName)
	req2.SetPathValue("metricValue", "20.3")
	rr2 := httptest.NewRecorder()
	Handle(rr2, req2)

	if rr2.Code != http.StatusOK {
		t.Errorf("Second update failed with status %v", rr2.Code)
	}

	value2, err := storage.GetGauge(metricName)
	if err != nil {
		t.Fatalf("GetGauge() after second update = %v", err)
	}
	if value2 != 20.3 {
		t.Errorf("After second update, gauge = %v, want 20.3", value2)
	}
}

func TestGetGaugeValue(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		want    float64
		wantErr bool
	}{
		{
			name:    "valid path with value",
			path:    "/10.5",
			want:    10.5,
			wantErr: false,
		},
		{
			name:    "valid path with integer",
			path:    "/100",
			want:    100.0,
			wantErr: false,
		},
		{
			name:    "path without value",
			path:    "/",
			want:    0,
			wantErr: true,
		},
		{
			name:    "empty path",
			path:    "",
			want:    0,
			wantErr: true,
		},
		{
			name:    "invalid value",
			path:    "/invalid",
			want:    0,
			wantErr: true,
		},
		{
			name:    "negative value",
			path:    "/-5.3",
			want:    -5.3,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getGaugeValue(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("getGaugeValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("getGaugeValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

