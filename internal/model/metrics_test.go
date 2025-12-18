package models

import (
	"testing"
)

func TestMetrics_Constants(t *testing.T) {
	if Counter != "counter" {
		t.Errorf("Expected Counter constant to be 'counter', got %q", Counter)
	}

	if Gauge != "gauge" {
		t.Errorf("Expected Gauge constant to be 'gauge', got %q", Gauge)
	}
}

func TestMetrics_Struct(t *testing.T) {
	// Тест создания структуры с counter
	delta := int64(10)
	metric := Metrics{
		ID:    "testCounter",
		MType: Counter,
		Delta: &delta,
	}

	if metric.ID != "testCounter" {
		t.Errorf("Expected ID 'testCounter', got %q", metric.ID)
	}

	if metric.MType != Counter {
		t.Errorf("Expected MType %q, got %q", Counter, metric.MType)
	}

	if metric.Delta == nil || *metric.Delta != 10 {
		t.Errorf("Expected Delta 10, got %v", metric.Delta)
	}

	if metric.Value != nil {
		t.Error("Expected Value to be nil for counter metric")
	}

	// Тест создания структуры с gauge
	value := 3.14
	metric2 := Metrics{
		ID:    "testGauge",
		MType: Gauge,
		Value: &value,
	}

	if metric2.ID != "testGauge" {
		t.Errorf("Expected ID 'testGauge', got %q", metric2.ID)
	}

	if metric2.MType != Gauge {
		t.Errorf("Expected MType %q, got %q", Gauge, metric2.MType)
	}

	if metric2.Value == nil || *metric2.Value != 3.14 {
		t.Errorf("Expected Value 3.14, got %v", metric2.Value)
	}

	if metric2.Delta != nil {
		t.Error("Expected Delta to be nil for gauge metric")
	}

	// Тест с нулевыми значениями
	zeroDelta := int64(0)
	zeroValue := 0.0
	metric3 := Metrics{
		Delta: &zeroDelta,
		Value: &zeroValue,
	}

	if metric3.Delta == nil || *metric3.Delta != 0 {
		t.Error("Zero delta should be preserved")
	}

	if metric3.Value == nil || *metric3.Value != 0.0 {
		t.Error("Zero value should be preserved")
	}
}

func TestMetrics_WithHash(t *testing.T) {
	metric := Metrics{
		Hash: "testHash123",
	}

	if metric.Hash != "testHash123" {
		t.Errorf("Expected Hash 'testHash123', got %q", metric.Hash)
	}
}
