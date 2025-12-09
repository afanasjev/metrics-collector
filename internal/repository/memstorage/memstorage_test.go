package memstorage

import (
	"sync"
	"testing"
)

func TestGetMemStorage(t *testing.T) {
	storage1 := GetMemStorage()
	storage2 := GetMemStorage()

	if storage1 != storage2 {
		t.Error("GetMemStorage should return the same instance (singleton)")
	}
}

func TestSetCounter(t *testing.T) {
	storage := GetMemStorage()

	tests := []struct {
		name     string
		metric   string
		value    int64
		expected int64
	}{
		{"set new counter", "testCounter", 10, 10},
		{"increment counter", "testCounter", 5, 15},
		{"increment again", "testCounter", 3, 18},
		{"new counter", "anotherCounter", 100, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := storage.SetCounter(tt.metric, tt.value)
			if err != nil {
				t.Errorf("SetCounter() error = %v", err)
			}

			got, err := storage.GetCounter(tt.metric)
			if err != nil {
				t.Fatalf("GetCounter() error = %v", err)
			}
			if got != tt.expected {
				t.Errorf("GetCounter() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestSetCounterConcurrent(t *testing.T) {
	storage := GetMemStorage()
	metricName := "concurrentCounter"
	goroutines := 100
	iterations := 10

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				_ = storage.SetCounter(metricName, 1)
			}
		}()
	}

	wg.Wait()

	got, err := storage.GetCounter(metricName)
	if err != nil {
		t.Fatalf("GetCounter() error after concurrent writes = %v", err)
	}
	expected := int64(goroutines * iterations)
	if got != expected {
		t.Errorf("GetCounter() = %v, want %v", got, expected)
	}
}

func TestSetGauge(t *testing.T) {
	storage := GetMemStorage()

	tests := []struct {
		name     string
		metric   string
		value    float64
		expected float64
	}{
		{"set new gauge", "testGauge", 10.5, 10.5},
		{"overwrite gauge", "testGauge", 20.3, 20.3},
		{"overwrite again", "testGauge", 15.7, 15.7},
		{"new gauge", "anotherGauge", 100.1, 100.1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := storage.SetGauge(tt.metric, tt.value)
			if err != nil {
				t.Errorf("SetGauge() error = %v", err)
			}

			got, err := storage.GetGauge(tt.metric)
			if err != nil {
				t.Fatalf("GetGauge() error = %v", err)
			}
			if got != tt.expected {
				t.Errorf("GetGauge() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestSetGaugeConcurrent(t *testing.T) {
	storage := GetMemStorage()
	metricName := "concurrentGauge"
	goroutines := 50

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func(val float64) {
			defer wg.Done()
			_ = storage.SetGauge(metricName, val)
		}(float64(i))
	}

	wg.Wait()

	got, err := storage.GetGauge(metricName)
	if err != nil {
		t.Fatalf("GetGauge() error after concurrent writes = %v", err)
	}
	// Последнее значение может быть любым из записанных
	if got < 0 || got >= float64(goroutines) {
		t.Errorf("GetGauge() = %v, should be in range [0, %v)", got, goroutines)
	}
}

func TestGetCounterNonExistent(t *testing.T) {
	storage := GetMemStorage()
	_, err := storage.GetCounter("nonExistent")
	if err == nil {
		t.Error("GetCounter() should return error for non-existent metric")
	}
}

func TestGetGaugeNonExistent(t *testing.T) {
	storage := GetMemStorage()
	_, err := storage.GetGauge("nonExistent")
	if err == nil {
		t.Error("GetGauge() should return error for non-existent metric")
	}
}

