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
		{name: "simple", path: "requests/123", want: "requests"},
		{name: "without metric value", path: "requests", want: "requests"},
		{name: "empty", path: "", wantErr: true},
		{name: "leading slash", path: "/requests", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetMetricName(tt.path)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetMetricName(%q) error = %v, wantErr=%v", tt.path, err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Fatalf("GetMetricName(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}

