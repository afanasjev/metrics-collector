package gauge

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/afanasjev/metrics-collector/internal/repository/memstorage"
)

func Set(w http.ResponseWriter, r *http.Request) {
	log.SetPrefix("[handler/counter/ServeHTTP")
	metricName := r.PathValue("metricName")

	metricValue, err := strconv.ParseFloat(r.PathValue("metricValue"), 64)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = SetGauge(metricName, metricValue)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func SetGauge(name string, value float64) error {
	storage := memstorage.GetMemStorage()
	return storage.SetGauge(name, value)

}

func GetGauge(name string) (float64, error) {
	storage := memstorage.GetMemStorage()
	return storage.GetGauge(name)

}

func Get(w http.ResponseWriter, r *http.Request) {
	log.SetPrefix("[handler/gauge/Get")
	metricName := r.PathValue("metricName")
	storage := memstorage.GetMemStorage()
	value, err := storage.GetGauge(metricName)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = fmt.Fprintf(w, "%v", value)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
