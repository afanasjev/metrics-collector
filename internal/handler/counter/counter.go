package counter

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

	metricValue, err := strconv.ParseInt(r.PathValue("metricValue"), 10, 64)
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

func Get(w http.ResponseWriter, r *http.Request) {
	log.SetPrefix("[handler/counter/Get")
	metricName := r.PathValue("metricName")
	storage := memstorage.GetMemStorage()
	value, err := storage.GetCounter(metricName)
	if err != nil {
		log.Print(err)
	}
	fmt.Fprintf(w, "%v", value)
}
