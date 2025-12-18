package memstorage

import (
	"fmt"
	"sync/atomic"
	"testing"
)

var uniqueID uint64

func uniqueMetricName(prefix string) string {
	id := atomic.AddUint64(&uniqueID, 1)
	return fmt.Sprintf("%s-%d", prefix, id)
}

func TestSetAndGetCounter(t *testing.T) {
	storage := GetMemStorage()
	name := uniqueMetricName("memstorage-counter")

	if err := storage.SetCounter(name, 5); err != nil {
		t.Fatalf("unexpected error from SetCounter: %v", err)
	}
	if err := storage.SetCounter(name, 3); err != nil {
		t.Fatalf("unexpected error from SetCounter: %v", err)
	}

	value, err := storage.GetCounter(name)
	if err != nil {
		t.Fatalf("unexpected error from GetCounter(%q): %v", name, err)
	}
	if want := int64(8); value != want {
		t.Fatalf("GetCounter(%q) = %d, want %d", name, value, want)
	}
}

func TestGetCounterMissing(t *testing.T) {
	storage := GetMemStorage()
	name := uniqueMetricName("memstorage-missing-counter")

	if _, err := storage.GetCounter(name); err == nil {
		t.Fatalf("expected error for missing counter %q", name)
	}
}

func TestSetAndGetGauge(t *testing.T) {
	storage := GetMemStorage()
	name := uniqueMetricName("memstorage-gauge")
	const expected = 3.14

	if err := storage.SetGauge(name, expected); err != nil {
		t.Fatalf("unexpected error from SetGauge: %v", err)
	}

	value, err := storage.GetGauge(name)
	if err != nil {
		t.Fatalf("unexpected error from GetGauge(%q): %v", name, err)
	}
	if value != expected {
		t.Fatalf("GetGauge(%q) = %v, want %v", name, value, expected)
	}
}

