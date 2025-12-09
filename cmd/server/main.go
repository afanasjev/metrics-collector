package main

import (
	"net/http"

	"github.com/afanasjev/metrics-collector/internal/handler/counter"
	"github.com/afanasjev/metrics-collector/internal/handler/gauge"

	"github.com/go-chi/chi/v5"
)

const (
	UpdatePath        = "/update/"
	CounterUpdatePath = "/update/counter/"
	GaugeUpdatePath   = "/update/gauge/"
)

func main() {
	run()
}

func run() {
	router := chi.NewRouter()
	router.Post("/update", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	router.Post("/update/{metricType}/{metricName}/{metricValue}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	})
	router.Post("/update/counter/{metricName}/{metricValue}", counter.Handle)
	router.Post("/update/gauge/{metricName}/{metricValue}", gauge.Handle)

	if err := http.ListenAndServe(":8080", router); err != nil {
		panic(err)
	}
}
