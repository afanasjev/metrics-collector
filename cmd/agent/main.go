package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"reflect"
	"runtime"
	"sync"
	"time"
)

const (
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
	serviceURL     = `http://localhost:8080`
	PollCount      = `PollCount`
	RandomValue    = `RandomValue`
)

func main() {
	var stats runtime.MemStats
	ctx := context.Background()
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go updateMetrics(ctx, wg, &stats)
	go sendMetrics(ctx, wg, &stats)

	wg.Wait()
}

func updateMetrics(ctx context.Context, wg *sync.WaitGroup, stats *runtime.MemStats) {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Print("shutting down updateMetrics...\n")
			wg.Done()
		case <-ticker.C:
			runtime.ReadMemStats(stats)
		}
	}
}

func sendMetrics(ctx context.Context, wg *sync.WaitGroup, stats *runtime.MemStats) {
	ticker := time.NewTicker(reportInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			wg.Done()
			fmt.Print("shutting down sendMetrics...\n")
		case <-ticker.C:
			sendGroup := &sync.WaitGroup{}
			sendGroup.Add(2)
			go sendGaugeUpdateRandom(ctx, sendGroup)
			go sendUpdateCounter(ctx, sendGroup)
			for _, fieldName := range runtimeMemStatsMetricsList() {
				v := reflect.ValueOf(stats).Elem()
				field := v.FieldByName(fieldName)
				if field.IsValid() {
					sendGroup.Add(2)
					go sendGaugeMetric(ctx, sendGroup, fieldName, fmt.Sprint(field))
					go sendUpdateCounter(ctx, sendGroup)
				} else {
					fmt.Printf("%s: INVALID\n", fieldName)
				}
			}
		}
	}
}

func runtimeMemStatsMetricsList() []string {
	return []string{"Alloc", "BuckHashSys", "Frees", "GCCPUFraction",
		"GCSys", "HeapAlloc", "HeapIdle", "HeapInuse", "HeapObjects", "HeapReleased",
		"HeapSys", "LastGC", "Lookups", "MCacheInuse", "MCacheSys", "MSpanInuse",
		"MSpanSys", "Mallocs", "NextGC", "NumForcedGC", "NumGC", "OtherSys",
		"PauseTotalNs", "StackInuse", "StackSys", "Sys", "TotalAlloc"}
}

func sendGaugeMetric(_ context.Context, wg *sync.WaitGroup, name string, value string) {
	metricURL := fmt.Sprintf("%s/update/gauge/%v/%v", serviceURL, name, value)
	response, err := http.Post(metricURL, "text/plain", nil)
	if err != nil {
		log.Printf("error sending gauge metric %v: %v\n", name, err)
	}

	defer func() {
		err := response.Body.Close()
		if err != nil {
			log.Printf("error on closing response body: %v", err)
		}
	}()
	wg.Done()
}

func sendGaugeUpdateRandom(_ context.Context, wg *sync.WaitGroup) {
	metricURL := fmt.Sprintf("%s/update/gauge/%v/%v", serviceURL, RandomValue, rand.Float64())
	response, err := http.Post(metricURL, "text/plain", nil)
	if err != nil {
		log.Printf("error sending gauge metric %v: %v\n", RandomValue, err)
	}

	defer func() {
		err := response.Body.Close()
		if err != nil {
			log.Printf("error on closing response body: %v", err)
		}
	}()

	wg.Done()
}

func sendUpdateCounter(_ context.Context, wg *sync.WaitGroup) {
	response, err := http.Post(
		fmt.Sprintf("%s/update/counter/%v/%v", serviceURL, PollCount, "1"),
		`text/plain`,
		nil)

	if err != nil {
		log.Printf("error sending PollCount metric: %v\n", err)
	}

	defer func() {
		err := response.Body.Close()
		if err != nil {
			log.Printf("error on closing response body: %v", err)
		}
	}()
	wg.Done()
}
