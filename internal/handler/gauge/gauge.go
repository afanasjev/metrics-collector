package gauge

import (
	"errors"
	"fmt"
	"github.com/afanasjev/metrics-collector/internal/handler"
	"github.com/afanasjev/metrics-collector/internal/repository/memstorage"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct{}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.SetPrefix("[handler/counter/ServeHTTP")
	if r.Method != http.MethodPost {
		log.Print("invalid http method")
		w.WriteHeader(http.StatusNotImplemented)
		return
	}
	metricName, err := handler.GetMetricName(r.URL.Path)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	metricValue, err := getGaugeValue(r.URL.Path)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	storage := memstorage.GetMemStorage()
	err = storage.SetGauge(metricName, metricValue)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Print(storage.PrintGauge())
}

func getGaugeValue(path string) (float64, error) {
	params := strings.Split(path, "/")
	if len(params) < 2 {
		return 0, errors.New("there is no metric value")
	}

	value, err := strconv.ParseFloat(params[1], 64)
	if err != nil {
		return 0, fmt.Errorf("invalid metric value: %s", params[1])
	}

	return value, nil
}
