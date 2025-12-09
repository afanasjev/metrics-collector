package counter

import (
	"log"
	"net/http"
	"strconv"

	"github.com/afanasjev/metrics-collector/internal/repository/memstorage"
)

type Handler struct{}

func Handle(w http.ResponseWriter, r *http.Request) {
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
