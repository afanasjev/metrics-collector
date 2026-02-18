package rest

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/afanasjev/metrics-collector/internal/config/server"

	"github.com/afanasjev/metrics-collector/internal/handler/counter"
	"github.com/afanasjev/metrics-collector/internal/handler/gauge"
	models "github.com/afanasjev/metrics-collector/internal/model"
)

func Set(w http.ResponseWriter, r *http.Request) {
	logger := server.GetServerConfig().GetLogger()
	data, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Errorf("Error reading body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	requestMetric := &models.Metrics{}
	responseMetric := &models.Metrics{}
	err = json.Unmarshal(data, requestMetric)
	if err != nil {
		logger.Errorf("Error parsing body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if requestMetric.MType == "counter" {
		responseMetric, err = setCounter(requestMetric.ID, *requestMetric.Delta)
		if err != nil {
			logger.Errorf("Error setting counter: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		logger.Debugf("Counter was set: %#v", responseMetric)
	}

	if requestMetric.MType == "gauge" {
		responseMetric, err = setGauge(requestMetric.ID, *requestMetric.Value)
		if err != nil {
			logger.Errorf("Error setting gauge: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		logger.Debugf("Gauge was set: %#v", responseMetric)
	}

	data, err = json.Marshal(responseMetric)
	if err != nil {
		logger.Errorf("error marshalling %#v: %v\n", responseMetric, err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)
	if err != nil {
		logger.Errorf("Error writing response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func Get(w http.ResponseWriter, r *http.Request) {
	logger := server.GetServerConfig().GetLogger()
	data, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Errorf("Error reading body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	requestMetric := &models.Metrics{}
	responseMetric := &models.Metrics{}
	err = json.Unmarshal(data, requestMetric)
	if err != nil {
		logger.Errorf("Error parsing body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if requestMetric.MType == "counter" {
		responseMetric, err = getCounter(requestMetric.ID)
		if err != nil {
			logger.Errorf("Error getting counter: %v", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}

	if requestMetric.MType == "gauge" {
		responseMetric, err = getGauge(requestMetric.ID)
		if err != nil {
			logger.Errorf("Error getting gauge: %v", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}

	data, err = json.Marshal(responseMetric)
	if err != nil {
		logger.Errorf("Error parsing body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)
	if err != nil {
		logger.Errorf("Error writing response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func setGauge(name string, value float64) (*models.Metrics, error) {
	err := gauge.SetGauge(name, value)
	if err != nil {
		return nil, err
	}

	return getGauge(name)
}

func getGauge(name string) (*models.Metrics, error) {
	value, err := gauge.GetGauge(name)
	if err != nil {
		return nil, err
	}

	return &models.Metrics{
		ID:    name,
		MType: "gauge",
		Value: &value,
	}, nil
}

func setCounter(name string, value int64) (*models.Metrics, error) {
	err := counter.SetCounter(name, value)
	if err != nil {
		return nil, err
	}

	return getCounter(name)
}

func getCounter(name string) (*models.Metrics, error) {
	value, err := counter.GetCounter(name)
	if err != nil {
		return nil, err
	}

	return &models.Metrics{
		ID:    name,
		MType: "counter",
		Delta: &value,
	}, nil
}
