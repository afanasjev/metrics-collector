package counter

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/afanasjev/metrics-collector/internal/repository/memstorage"
)

func Set(w http.ResponseWriter, r *http.Request) {
	log.SetPrefix("[handler/counter/ServeHTTP")
	metricValue, err := strconv.ParseInt(r.PathValue("metricValue"), 10, 64)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = SetCounter(r.PathValue("metricName"), metricValue)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func Get(w http.ResponseWriter, r *http.Request) {
	log.SetPrefix("[handler/counter/Get")
	metricName := r.PathValue("metricName")
	storage := memstorage.GetMemStorage()
	value, err := storage.GetCounter(metricName)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	_, err = fmt.Fprintf(w, "%v", value)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func SetCounter(name string, value int64) error {
	storage := memstorage.GetMemStorage()
	return storage.SetCounter(name, value)
}
func GetCounter(name string) (int64, error) {
	storage := memstorage.GetMemStorage()
	return storage.GetCounter(name)
}
