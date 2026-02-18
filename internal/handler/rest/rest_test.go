package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/afanasjev/metrics-collector/internal/handler/testutil"
	models "github.com/afanasjev/metrics-collector/internal/model"
	"github.com/afanasjev/metrics-collector/internal/repository/memstorage"
)

func TestSetCounterJSON(t *testing.T) {
	storage := memstorage.GetMemStorage()
	metricName := testutil.MetricName(t, "rest-counter-set")
	payload := fmt.Sprintf(`{"id":"%s","type":"counter","delta":%d}`, metricName, 7)

	req := httptest.NewRequest(http.MethodPost, "/update", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	Set(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d; want %d", rr.Code, http.StatusOK)
	}
	if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("Content-Type = %q, want %q", ct, "application/json")
	}

	var response models.Metrics
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("unexpected body: %v", err)
	}
	if response.ID != metricName || response.MType != "counter" {
		t.Fatalf("response metric = %#v", response)
	}
	if response.Delta == nil || *response.Delta != 7 {
		t.Fatalf("delta = %v, want %d", response.Delta, 7)
	}

	value, err := storage.GetCounter(metricName)
	if err != nil {
		t.Fatalf("GetCounter failed: %v", err)
	}
	if value != 7 {
		t.Fatalf("counter %s = %d, want %d", metricName, value, 7)
	}
}

func TestSetGaugeJSON(t *testing.T) {
	storage := memstorage.GetMemStorage()
	metricName := testutil.MetricName(t, "rest-gauge-set")
	payload := fmt.Sprintf(`{"id":"%s","type":"gauge","value":%g}`, metricName, 3.21)

	req := httptest.NewRequest(http.MethodPost, "/update", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	Set(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d; want %d", rr.Code, http.StatusOK)
	}
	if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("Content-Type = %q, want %q", ct, "application/json")
	}

	var response models.Metrics
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("unexpected body: %v", err)
	}
	if response.ID != metricName || response.MType != "gauge" {
		t.Fatalf("response metric = %#v", response)
	}
	if response.Value == nil || *response.Value != 3.21 {
		t.Fatalf("value = %v, want %v", response.Value, 3.21)
	}

	value, err := storage.GetGauge(metricName)
	if err != nil {
		t.Fatalf("GetGauge failed: %v", err)
	}
	if value != 3.21 {
		t.Fatalf("gauge %s = %v, want %v", metricName, value, 3.21)
	}
}

func TestSetInvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/update", strings.NewReader("{invalid json"))
	rr := httptest.NewRecorder()

	Set(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d; want %d", rr.Code, http.StatusBadRequest)
	}
}

func TestSetReadError(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/update", nil)
	req.Body = errReadCloser{err: errors.New("boom")}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	Set(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d; want %d", rr.Code, http.StatusInternalServerError)
	}
}

func TestSetWriteFailure(t *testing.T) {
	metricName := testutil.MetricName(t, "rest-set-write")
	payload := fmt.Sprintf(`{"id":"%s","type":"gauge","value":%g}`, metricName, 5.5)
	req := httptest.NewRequest(http.MethodPost, "/update", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	writer := testutil.NewFailingResponseWriter(errors.New("boom"))

	Set(writer, req)

	if writer.Status() != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", writer.Status(), http.StatusInternalServerError)
	}
}

func TestGetJSON(t *testing.T) {
	storage := memstorage.GetMemStorage()
	metricName := testutil.MetricName(t, "rest-get")
	if err := storage.SetCounter(metricName, 13); err != nil {
		t.Fatalf("SetCounter failed: %v", err)
	}

	payload := fmt.Sprintf(`{"id":"%s","type":"counter"}`, metricName)
	req := httptest.NewRequest(http.MethodPost, "/value", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	Get(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d; want %d", rr.Code, http.StatusOK)
	}
	if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("Content-Type = %q, want %q", ct, "application/json")
	}

	var resp models.Metrics
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unexpected body: %v", err)
	}
	if resp.ID != metricName {
		t.Fatalf("metric ID = %q, want %q", resp.ID, metricName)
	}
	if resp.Delta == nil || *resp.Delta != 13 {
		t.Fatalf("delta = %v, want %d", resp.Delta, 13)
	}
}

func TestGetInvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/value", strings.NewReader("{invalid json"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	Get(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
	}
}

func TestGetWriteFailure(t *testing.T) {
	storage := memstorage.GetMemStorage()
	metricName := testutil.MetricName(t, "rest-get-write")
	if err := storage.SetCounter(metricName, 5); err != nil {
		t.Fatalf("SetCounter failed: %v", err)
	}

	payload := fmt.Sprintf(`{"id":"%s","type":"counter"}`, metricName)
	req := httptest.NewRequest(http.MethodPost, "/value", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	writer := testutil.NewFailingResponseWriter(errors.New("boom"))

	Get(writer, req)

	if writer.Status() != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", writer.Status(), http.StatusInternalServerError)
	}
}

type errReadCloser struct {
	err error
}

func (e errReadCloser) Read([]byte) (int, error) {
	return 0, e.err
}

func (e errReadCloser) Close() error {
	return nil
}
