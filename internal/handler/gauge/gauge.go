package gauge

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/afanasjev/metrics-collector/internal/repository/memstorage"
)

type Handler struct{}

func Set(w http.ResponseWriter, r *http.Request) {
	log.SetPrefix("[handler/counter/ServeHTTP")
	metricName := r.PathValue("metricName")

	metricValue, err := strconv.ParseFloat(r.PathValue("metricValue"), 64)
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

func Get(w http.ResponseWriter, r *http.Request) {
	log.SetPrefix("[handler/gauge/Get")
	metricName := r.PathValue("metricName")
	storage := memstorage.GetMemStorage()
	value, err := storage.GetGauge(metricName)
	if err != nil {
		log.Print(err)
	}
	fmt.Fprintf(w, "%v", value)
}
