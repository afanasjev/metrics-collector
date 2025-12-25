package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/afanasjev/metrics-collector/internal/handler/counter"
	"github.com/afanasjev/metrics-collector/internal/handler/gauge"
	"github.com/caarlos0/env/v11"
	"github.com/go-chi/chi/v5"
)

type Configuration struct {
	ServerAddress string `env:"ADDRESS"`
}

var cfg Configuration

func main() {

	cfg = Configuration{}
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	if cfg.ServerAddress == "" {
		flag.StringVar(&cfg.ServerAddress, "a", "localhost:8080", "address to listen on")
		flag.Parse()
	}

	fmt.Printf("Configuration: %+v\n", cfg)

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
	router.Post("/update/counter/{metricName}/{metricValue}", counter.Set)
	router.Post("/update/gauge/{metricName}/{metricValue}", gauge.Set)
	router.Get("/value/counter/{metricName}", counter.Get)
	router.Get("/value/gauge/{metricName}", gauge.Get)

	if err := http.ListenAndServe(cfg.ServerAddress, router); err != nil {
		log.Fatal(err)
	}
}
