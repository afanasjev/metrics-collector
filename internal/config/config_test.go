package config

import (
	"testing"

	"github.com/afanasjev/metrics-collector/internal/config/agent"
	"github.com/afanasjev/metrics-collector/internal/config/server"
	"go.uber.org/zap"
)

func TestAgentConfigurationGetters(t *testing.T) {
	logger := zap.NewNop()
	cfg := &agent.Configuration{
		ServerAddress: "http://localhost:8080",
		Logger:        logger,
	}

	if got := cfg.GetServerAddress(); got != "http://localhost:8080" {
		t.Fatalf("GetServerAddress = %q, want %q", got, "http://localhost:8080")
	}
	if cfg.GetLogger() == nil {
		t.Fatalf("GetLogger returned nil")
	}
}

func TestAgentConfigurationFallbackLogger(t *testing.T) {
	cfg := &agent.Configuration{}
	if cfg.GetLogger() == nil {
		t.Fatalf("GetLogger should return fallback logger for nil Logger")
	}
}

func TestServerConfigurationGetters(t *testing.T) {
	logger := zap.NewNop()
	cfg := &server.Configuration{
		ServerAddress: "localhost:8080",
		Debug:         true,
		Logger:        logger,
	}

	if got := cfg.GetServerAddress(); got != "localhost:8080" {
		t.Fatalf("GetServerAddress = %q, want %q", got, "localhost:8080")
	}
	if cfg.GetLogger() == nil {
		t.Fatalf("GetLogger returned nil")
	}
}

func TestServerConfigurationFallbackLogger(t *testing.T) {
	cfg := &server.Configuration{}
	if cfg.GetLogger() == nil {
		t.Fatalf("GetLogger should return fallback logger for nil Logger")
	}
}
