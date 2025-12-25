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

func TestCounterStorage(t *testing.T) {
	storage := GetMemStorage()

	t.Run("accumulates values", func(t *testing.T) {
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
	})

	t.Run("missing metric returns error", func(t *testing.T) {
		if _, err := storage.GetCounter(uniqueMetricName("memstorage-missing-counter")); err == nil {
			t.Fatalf("expected error for missing counter")
		}
	})
}

func TestGaugeStorage(t *testing.T) {
	storage := GetMemStorage()

	t.Run("stores value", func(t *testing.T) {
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
	})

	t.Run("missing metric returns error", func(t *testing.T) {
		if _, err := storage.GetGauge(uniqueMetricName("memstorage-missing-gauge")); err == nil {
			t.Fatalf("expected error for missing gauge")
		}
	})
}
