package main

import (
	"net/http"
	"time"

	"github.com/afanasjev/metrics-collector/internal/handler/counter"
	"github.com/afanasjev/metrics-collector/internal/handler/gauge"
	"github.com/afanasjev/metrics-collector/internal/server/config"
	"github.com/go-chi/chi/v5"
)

type responseData struct {
	statusCode int
	size       int
}

type loggingWithResponseData struct {
	http.ResponseWriter
	responseData *responseData
}

func (lwrd *loggingWithResponseData) Write(data []byte) (int, error) {
	responseSize, err := lwrd.ResponseWriter.Write(data)
	if err != nil {
		return 0, err
	}
	lwrd.responseData.size += responseSize
	return responseSize, nil
}

func (lwrd *loggingWithResponseData) WriteHeader(statusCode int) {
	lwrd.ResponseWriter.WriteHeader(statusCode)
	lwrd.responseData.statusCode = statusCode
}

func WithLogging(h http.Handler) http.Handler {
	wrapped := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rd := &responseData{}
		lwrd := &loggingWithResponseData{
			ResponseWriter: w,
			responseData:   rd,
		}
		h.ServeHTTP(lwrd, r)

		sugar := config.Get().GetLogger()

		sugar.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"status", lwrd.responseData.statusCode,
			"size", lwrd.responseData.size,
			"duration", time.Since(start),
		)

	}
	return http.HandlerFunc(wrapped)
}

func main() {

	config.NewConfig()
	run()
}

func run() {
	router := chi.NewRouter()
	router.Use(WithLogging)
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

	cfg := config.Get()
	logger := cfg.GetLogger()
	logger.Infow("Starting server on " + cfg.GetServerAddress())
	if err := http.ListenAndServe(cfg.GetServerAddress(), router); err != nil {
		logger.Fatalw(err.Error(), "event", "start server")
	}
}
