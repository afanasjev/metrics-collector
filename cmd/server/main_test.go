package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/afanasjev/metrics-collector/internal/config/server"
	"go.uber.org/zap"
)

func initServerConfigForTests(t *testing.T) {
	t.Helper()
	cfg = &server.Configuration{
		ServerAddress: "127.0.0.1:8092",
		Debug:         true,
		Logger:        zap.NewNop(),
	}
	lggr = cfg.GetLogger()
}

func TestLoggingWithResponseDataWriteHeader(t *testing.T) {
	recorder := httptest.NewRecorder()
	rd := &responseData{}
	w := &loggingWithResponseData{
		w:            recorder,
		responseData: rd,
	}

	w.WriteHeader(http.StatusCreated)
	if rd.statusCode != http.StatusCreated {
		t.Fatalf("statusCode = %d, want %d", rd.statusCode, http.StatusCreated)
	}
	if recorder.Code != http.StatusCreated {
		t.Fatalf("recorder.Code = %d, want %d", recorder.Code, http.StatusCreated)
	}
}

func TestLoggingWithResponseDataWriteSetsDefaultStatusAndSize(t *testing.T) {
	recorder := httptest.NewRecorder()
	rd := &responseData{}
	w := &loggingWithResponseData{
		w:            recorder,
		responseData: rd,
	}

	n, err := w.Write([]byte("abc"))
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}
	if n != 3 {
		t.Fatalf("Write bytes = %d, want %d", n, 3)
	}
	if rd.statusCode != http.StatusOK {
		t.Fatalf("default status = %d, want %d", rd.statusCode, http.StatusOK)
	}
	if rd.size != 3 {
		t.Fatalf("response size = %d, want %d", rd.size, 3)
	}
}

func TestWithLoggingPassesResponseThrough(t *testing.T) {
	initServerConfigForTests(t)

	handler := WithLogging(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		_, _ = w.Write([]byte("ok"))
	}))

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusAccepted {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusAccepted)
	}
	if rr.Body.String() != "ok" {
		t.Fatalf("body = %q, want %q", rr.Body.String(), "ok")
	}
}

func TestNormalizeListenAddress(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "localhost", input: "localhost:8080", want: ":8080"},
		{name: "ipv4 loopback", input: "127.0.0.1:8081", want: ":8081"},
		{name: "ipv6 loopback", input: "[::1]:8082", want: ":8082"},
		{name: "custom host unchanged", input: "0.0.0.0:8083", want: "0.0.0.0:8083"},
		{name: "invalid unchanged", input: "8084", want: "8084"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizeListenAddress(tt.input)
			if got != tt.want {
				t.Fatalf("normalizeListenAddress(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
