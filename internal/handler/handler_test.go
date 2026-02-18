package handler

import "testing"

func TestGetMetricName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		path    string
		want    string
		wantErr bool
	}{
		{name: "simple", path: "requests/123", want: "requests"},
		{name: "single segment", path: "requests", want: "requests"},
		{name: "multiple hops", path: "cpu/usage/percent", want: "cpu"},
		{name: "empty path", path: "", wantErr: true},
		{name: "leading slash", path: "/requests", wantErr: true},
		{name: "extra slash", path: "mem//usage", want: "mem"},
		{name: "root slash", path: "/", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetMetricName(tt.path)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error for %q", tt.path)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error for %q: %v", tt.path, err)
			}
			if got != tt.want {
				t.Fatalf("GetMetricName(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}
