package main

import (
	"net"
	"net/http"
	"time"

	"github.com/afanasjev/metrics-collector/internal/config/server"
	"github.com/afanasjev/metrics-collector/internal/handler/counter"
	"github.com/afanasjev/metrics-collector/internal/handler/gauge"
	"github.com/afanasjev/metrics-collector/internal/handler/rest"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"moul.io/http2curl"
)

var cfg *server.Configuration
var lggr *zap.SugaredLogger

type (
	responseData struct {
		statusCode int
		size       int
	}

	loggingWithResponseData struct {
		w            http.ResponseWriter // Не встраиваем, а храним как поле
		responseData *responseData
	}
)

func (lwrd *loggingWithResponseData) Header() http.Header {
	return lwrd.w.Header()
}

func (lwrd *loggingWithResponseData) Write(data []byte) (int, error) {
	if lwrd.responseData.statusCode == 0 {
		lwrd.WriteHeader(http.StatusOK)
	}
	responseSize, err := lwrd.w.Write(data)
	lwrd.responseData.size += responseSize
	return responseSize, err
}

func (lwrd *loggingWithResponseData) WriteHeader(statusCode int) {
	lwrd.w.WriteHeader(statusCode)
	lwrd.responseData.statusCode = statusCode
}

func WithLogging(next http.Handler) http.Handler {
	wrapped := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rd := &responseData{}
		lwrd := &loggingWithResponseData{
			w:            w,
			responseData: rd,
		}
		next.ServeHTTP(lwrd, r)

		lggr.Debug(http2curl.GetCurlCommand(r))
		lggr.Infoln(
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
	cfg = server.NewConfig()
	lggr = cfg.GetLogger()
	lggr.Infof("Configuration loaded: %#v\n", cfg)
	run()
}

func run() {
	router := chi.NewRouter()
	router.Use(WithLogging)

	router.Post("/update/", rest.Set)
	router.Post("/value/", rest.Get)

	router.Post("/update/{metricType}/{metricName}/{metricValue}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	})
	router.Post("/update/counter/{metricName}/{metricValue}", counter.Set)
	router.Post("/update/gauge/{metricName}/{metricValue}", gauge.Set)
	router.Get("/value/counter/{metricName}", counter.Get)
	router.Get("/value/gauge/{metricName}", gauge.Get)

	listenAddress := normalizeListenAddress(cfg.GetServerAddress())
	lggr.Infow("Starting server on " + listenAddress)
	if err := http.ListenAndServe(listenAddress, router); err != nil {
		lggr.Fatalw(err.Error(), "event", "start server")
	}
}

// normalizeListenAddress avoids IPv4/IPv6 localhost mismatch by listening on all
// interfaces for loopback addresses while preserving the configured port.
func normalizeListenAddress(address string) string {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return address
	}

	switch host {
	case "localhost", "127.0.0.1", "::1":
		return ":" + port
	default:
		return address
	}
}
