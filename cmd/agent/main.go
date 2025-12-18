package main

import (
	"context"
	"flag"
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
	pollCount   = `PollCount`
	randomValue = `RandomValue`
)

var ServerAddress string
var PollInterval int
var ReportInterval int

func main() {

	flag.StringVar(&ServerAddress, "a", "localhost:8080", "The server address in the format of host:port")
	flag.IntVar(&PollInterval, "p", 10, "The poll interval in seconds")
	flag.IntVar(&ReportInterval, "r", 10, "The report interval in seconds")
	flag.Parse()
	ServerAddress = fmt.Sprintf("http://%s", ServerAddress)

	var stats runtime.MemStats
	ctx := context.Background()
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go updateMetrics(ctx, wg, &stats)
	go sendMetrics(ctx, wg, &stats)

	wg.Wait()
}

func updateMetrics(ctx context.Context, wg *sync.WaitGroup, stats *runtime.MemStats) {
	ticker := time.NewTicker(time.Duration(PollInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Print("shutting down updateMetrics...\n")
			wg.Done()
		case <-ticker.C:
			runtime.ReadMemStats(stats)
			wg.Add(1)
			go sendUpdateCounter(ctx, wg, len(runtimeMemStatsMetricsList()))
		}
	}
}

func sendMetrics(ctx context.Context, wg *sync.WaitGroup, stats *runtime.MemStats) {
	ticker := time.NewTicker(time.Duration(ReportInterval) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			wg.Done()
			fmt.Print("shutting down sendMetrics...\n")
		case <-ticker.C:
			sendGroup := &sync.WaitGroup{}
			sendGroup.Add(1)
			go sendGaugeUpdateRandom(ctx, sendGroup)
			for _, fieldName := range runtimeMemStatsMetricsList() {
				v := reflect.ValueOf(stats).Elem()
				field := v.FieldByName(fieldName)
				if field.IsValid() {
					sendGroup.Add(1)
					go sendGaugeMetric(ctx, sendGroup, fieldName, fmt.Sprint(field))
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
	metricURL := fmt.Sprintf("%s/update/gauge/%v/%v", ServerAddress, name, value)
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
	metricURL := fmt.Sprintf("%s/update/gauge/%v/%v", ServerAddress, randomValue, rand.Float64())
	response, err := http.Post(metricURL, "text/plain", nil)
	if err != nil {
		log.Printf("error sending gauge metric %v: %v\n", randomValue, err)
	}

	defer func() {
		err := response.Body.Close()
		if err != nil {
			log.Printf("error on closing response body: %v", err)
		}
	}()

	wg.Done()
}

func sendUpdateCounter(_ context.Context, wg *sync.WaitGroup, counter int) {
	response, err := http.Post(
		fmt.Sprintf("%s/update/counter/%v/%v", ServerAddress, pollCount, counter),
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
