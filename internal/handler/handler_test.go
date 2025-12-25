package handler

import (
	"testing"
)

func TestGetMetricName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		path    string
		want    string
		wantErr bool
	}{
		{name: "simple", path: "requests/123", want: "requests"},
		{name: "only name", path: "requests", want: "requests"},
		{name: "with trailing slash", path: "requests/", want: "requests"},
		{name: "multi segment", path: "cpu/usage/percent", want: "cpu"},
		{name: "empty path", path: "", wantErr: true},
		{name: "leading slash", path: "/requests", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := GetMetricName(tt.path)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetMetricName(%q) error = %v, wantErr=%v", tt.path, err, tt.wantErr)
			}
			if err == nil && got != tt.want {
				t.Fatalf("GetMetricName(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}
