package counter

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

	metricValue, err := getCounterValue(r.URL.Path)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	storage := memstorage.GetMemStorage()
	err = storage.SetCounter(metricName, metricValue)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Print(storage.PrintCounter())
}

func getCounterValue(path string) (int64, error) {
	params := strings.Split(path, "/")
	if len(params) < 2 {
		return 0, errors.New("there is no metric value")
	}

	value, err := strconv.ParseInt(params[1], 10, 64)
	if err != nil || value < 0 {
		return 0, fmt.Errorf("invalid metric value: %s", params[1])
	}

	return value, nil
}
