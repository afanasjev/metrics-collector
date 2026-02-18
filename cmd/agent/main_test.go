package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/afanasjev/metrics-collector/internal/config/agent"
	models "github.com/afanasjev/metrics-collector/internal/model"
	"go.uber.org/zap"
)

func initAgentConfigForTests(t *testing.T, address string) {
	t.Helper()
	cfg = &agent.Configuration{
		ServerAddress:  address,
		PollInterval:   1,
		ReportInterval: 1,
		Logger:         zap.NewNop(),
	}
	logger = cfg.GetLogger()
}

func TestRuntimeMemStatsMetricsList(t *testing.T) {
	metrics := NewMetrics()
	if len(metrics.metrics) == 0 {
		t.Fatalf("metrics map should not be empty")
	}

	if _, exists := metrics.metrics["Alloc"]; !exists {
		t.Fatalf("Alloc should be present in runtime metrics map")
	}
}

func TestSendGaugeMetricSendsPayload(t *testing.T) {
	var got models.Metrics
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatalf("decode failed: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	initAgentConfigForTests(t, server.URL)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	sendGauge(wg, "HeapAlloc", 42.5)
	wg.Wait()

	if got.MType != "gauge" {
		t.Fatalf("type = %q, want %q", got.MType, "gauge")
	}
	if got.ID != "HeapAlloc" {
		t.Fatalf("id = %q, want %q", got.ID, "HeapAlloc")
	}
	if got.Value == nil || *got.Value != 42.5 {
		t.Fatalf("value = %v, want %v", got.Value, 42.5)
	}
}

func TestSendUpdateCounterSendsPayload(t *testing.T) {
	var got models.Metrics
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatalf("decode failed: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	initAgentConfigForTests(t, server.URL)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	counter := &PollCounter{value: 123}
	sendUpdateCounter(context.Background(), wg, counter)
	wg.Wait()

	if got.MType != "counter" {
		t.Fatalf("type = %q, want %q", got.MType, "counter")
	}
	if got.ID != pollCount {
		t.Fatalf("id = %q, want %q", got.ID, pollCount)
	}
	if got.Delta == nil || *got.Delta != 123 {
		t.Fatalf("delta = %v, want %d", got.Delta, 123)
	}
}
