package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"reflect"
	"runtime"
	"sync"
	"time"

	"github.com/afanasjev/metrics-collector/internal/config/agent"
	models "github.com/afanasjev/metrics-collector/internal/model"
	"go.uber.org/zap"
	"moul.io/http2curl"
)

const (
	pollCount       = `PollCount`
	randomValue     = `RandomValue`
	applicationJSON = `application/json`
)

type PollCounter struct {
	value int64
	mux   sync.Mutex
}

func (p *PollCounter) Inc() {
	p.mux.Lock()
	defer p.mux.Unlock()
	p.value++
}

type Metrics struct {
	metrics map[string]float64
	mux     sync.Mutex
}

var logger *zap.SugaredLogger
var cfg *agent.Configuration

func main() {
	cfg = agent.NewConfig()
	logger = cfg.GetLogger()
	logger.Infof("Server address: %v", cfg.GetServerAddress())
	run()
}

func run() {
	ctx := context.Background()
	wg := &sync.WaitGroup{}
	wg.Add(2)

	pollCnt := &PollCounter{}
	metrics := NewMetrics()
	go updateMetrics(ctx, wg, metrics, pollCnt)
	go sendMetrics(ctx, wg, metrics, pollCnt)

	wg.Wait()
}

func updateMetrics(ctx context.Context, wg *sync.WaitGroup,
	metrics *Metrics, pollCnt *PollCounter) {
	ticker := time.NewTicker(time.Duration(cfg.PollInterval) * time.Second)
	defer ticker.Stop()
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			logger.Info("shutting down updateMetrics...\n")
			wg.Done()
		case <-ticker.C:
			updateMetricsFromMemStats(metrics, pollCnt)
			logger.Info("Metrics updated.")
		}
	}
}

func sendMetrics(ctx context.Context, wg *sync.WaitGroup,
	metrics *Metrics, cnt *PollCounter) {
	ticker := time.NewTicker(time.Duration(cfg.ReportInterval) * time.Second)
	defer ticker.Stop()
	defer wg.Done()
	sendGaugeWg := &sync.WaitGroup{}
	defer sendGaugeWg.Wait()

	for {
		select {
		case <-ctx.Done():
			wg.Done()
			fmt.Print("shutting down sendMetrics...\n")
		case <-ticker.C:
			sendGaugeWg.Add(2)
			go sendGauge(sendGaugeWg, randomValue, rand.Float64())
			go sendUpdateCounter(ctx, sendGaugeWg, cnt)
			for metricName := range metrics.metrics {
				sendGaugeWg.Add(1)
				go sendGauge(sendGaugeWg, metricName, metrics.metrics[metricName])
			}
		}
	}
}

func sendGauge(wg *sync.WaitGroup, name string, value float64) {
	defer wg.Done()
	metricURL := fmt.Sprintf("%s/update/", cfg.GetServerAddress())
	metric := models.Metrics{
		ID:    name,
		Value: &value,
		MType: "gauge",
	}
	data, err := json.Marshal(metric)
	if err != nil {
		logger.Errorf("error marshalling %#v: %v\n", metric, err)
	}
	response, err := http.Post(metricURL, applicationJSON, bytes.NewReader(data))
	if err != nil {
		logger.Errorf("error sending gauge metric %v: %v\nerror: %v", name, string(data), err)
	}
	logger.Debugf("URL: %v DATA %v\n", metricURL, string(data))

	defer func() {
		err := response.Body.Close()
		if err != nil {
			logger.Errorf("error on closing response body: %v", err)
		}
	}()
}

func sendUpdateCounter(_ context.Context, wg *sync.WaitGroup, counter *PollCounter) {
	defer wg.Done()

	counter.mux.Lock()
	defer counter.mux.Unlock()

	metric := models.Metrics{
		ID:    pollCount,
		Delta: &counter.value,
		MType: "counter",
	}
	data, err := json.Marshal(metric)
	if err != nil {
		logger.Errorf("error marshalling %#v: %v\n", metric, err)
	}

	metricURL := fmt.Sprintf("%s/update", cfg.GetServerAddress())
	response, err := http.Post(metricURL, applicationJSON, bytes.NewReader(data))
	if err == nil {
		logger.Debug(http2curl.GetCurlCommand(response.Request))
	}

	if err != nil {
		logger.Errorf("error sending PollCount metric: %v\n", err)
	}

	logger.Debugf("URL: %v DATA %v\n", metricURL, string(data))

	defer func() {
		err := response.Body.Close()
		if err != nil {
			logger.Errorf("error on closing response body: %v", err)
		}
	}()
}

func updateMetricsFromMemStats(metrics *Metrics, pollCnt *PollCounter) {
	stats := &runtime.MemStats{}
	metrics.mux.Lock()
	defer metrics.mux.Unlock()
	runtime.ReadMemStats(stats)
	v := reflect.ValueOf(stats)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return
	}
	for fieldName := range metrics.metrics {
		field := v.FieldByName(fieldName)
		if !field.IsValid() {
			continue
		}

		switch field.Kind() {
		case reflect.Float32, reflect.Float64:
			v := field.Float()
			if v != metrics.metrics[fieldName] {
				pollCnt.Inc()
				metrics.metrics[fieldName] = v
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			v := float64(field.Uint())
			if v != metrics.metrics[fieldName] {
				pollCnt.Inc()
				metrics.metrics[fieldName] = v
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v := float64(field.Int())
			if v != metrics.metrics[fieldName] {
				pollCnt.Inc()
				metrics.metrics[fieldName] = v
			}
		}
	}
}
func NewMetrics() *Metrics {
	return &Metrics{
		metrics: map[string]float64{
			"Alloc":         0,
			"BuckHashSys":   0,
			"Frees":         0,
			"GCCPUFraction": 0,
			"GCSys":         0,
			"HeapAlloc":     0,
			"HeapIdle":      0,
			"HeapInuse":     0,
			"HeapObjects":   0,
			"HeapReleased":  0,
			"HeapSys":       0,
			"LastGC":        0,
			"Lookups":       0,
			"MCacheInuse":   0,
			"MCacheSys":     0,
			"MSpanInuse":    0,
			"MSpanSys":      0,
			"Mallocs":       0,
			"NextGC":        0,
			"NumForcedGC":   0,
			"NumGC":         0,
			"OtherSys":      0,
			"PauseTotalNs":  0,
			"StackInuse":    0,
			"StackSys":      0,
			"Sys":           0,
			"TotalAlloc":    0,
		},
	}
}
