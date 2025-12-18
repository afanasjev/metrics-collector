package main

import (
	"github.com/afanasjev/metrics-collector/internal/handler/counter"
	"github.com/afanasjev/metrics-collector/internal/handler/update"
	"net/http"
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
	mux := http.NewServeMux()
	mux.Handle(UpdatePath, http.StripPrefix(UpdatePath, &update.Handler{}))
	mux.Handle(CounterUpdatePath, http.StripPrefix(CounterUpdatePath, &counter.Handler{}))
	mux.Handle(GaugeUpdatePath, http.StripPrefix(GaugeUpdatePath, &counter.Handler{}))
	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}
