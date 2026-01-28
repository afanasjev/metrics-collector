package models

import (
	"encoding/json"
	"testing"
)

func TestMetricTypeConstants(t *testing.T) {
	if Counter != "counter" {
		t.Fatalf("Counter = %q, want %q", Counter, "counter")
	}
	if Gauge != "gauge" {
		t.Fatalf("Gauge = %q, want %q", Gauge, "gauge")
	}
}

func TestMetricsJSONEncodingOmitempty(t *testing.T) {
	value := 10.5
	m := Metrics{
		ID:    "cpu",
		MType: Gauge,
		Value: &value,
	}

	data, err := json.Marshal(m)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var decoded map[string]any
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if _, ok := decoded["delta"]; ok {
		t.Fatalf("delta should be omitted when nil")
	}
	if got, ok := decoded["id"].(string); !ok || got != "cpu" {
		t.Fatalf("id = %#v, want %q", decoded["id"], "cpu")
	}
}
