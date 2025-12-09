package handler

import (
	"testing"
)

func TestGetMetricName(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		want    string
		wantErr bool
	}{
		{
			name:    "valid path with single metric",
			path:    "testMetric",
			want:    "testMetric",
			wantErr: false,
		},
		{
			name:    "valid path with slash",
			path:    "/testMetric",
			want:    "",
			wantErr: true, // strings.Split("/testMetric", "/") returns ["", "testMetric"], so params[0] is ""
		},
		{
			name:    "empty path",
			path:    "",
			want:    "",
			wantErr: true,
		},
		{
			name:    "path starting with slash only",
			path:    "/",
			want:    "",
			wantErr: true,
		},
		{
			name:    "path with multiple slashes",
			path:    "/update/gauge/testMetric/123.45",
			want:    "",
			wantErr: true, // strings.Split returns ["", "update", "gauge", ...], so params[0] is ""
		},
		{
			name:    "metric name with underscores",
			path:    "test_metric_name",
			want:    "test_metric_name",
			wantErr: false,
		},
		{
			name:    "path without leading slash",
			path:    "update/gauge/testMetric",
			want:    "update",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetMetricName(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMetricName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetMetricName() = %v, want %v", got, tt.want)
			}
		})
	}
}

