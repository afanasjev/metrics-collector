package memstorage

import (
	"fmt"
	"sync"

	"github.com/afanasjev/metrics-collector/internal/config/server"
	"go.uber.org/zap"
)

type MemStorage struct {
	counter counter
	gauge   gauge
}

type counter struct {
	data  map[string]int64
	mutex sync.Mutex
}

type gauge struct {
	data  map[string]float64
	mutex sync.Mutex
}

var storage *MemStorage
var once sync.Once
var logger *zap.SugaredLogger
var cfg *server.Configuration

func GetMemStorage() *MemStorage {
	once.Do(func() {
		cfg = server.GetServerConfig()
		logger = cfg.GetLogger()
		storage = &MemStorage{
			counter: counter{
				data:  make(map[string]int64),
				mutex: sync.Mutex{},
			},
			gauge: gauge{
				data:  make(map[string]float64),
				mutex: sync.Mutex{},
			},
		}
	})
	return storage
}

func (ms *MemStorage) SetCounter(name string, value int64) error {
	ms.counter.mutex.Lock()
	defer ms.counter.mutex.Unlock()
	ms.counter.data[name] += value
	logger.Debugf("memstorage after Set counter: %#v", ms.counter.data)
	return nil
}

func (ms *MemStorage) SetGauge(name string, value float64) error {
	ms.gauge.mutex.Lock()
	defer ms.gauge.mutex.Unlock()
	ms.gauge.data[name] = value
	logger.Debugf("memstorage after Set gauge: %#v", ms.gauge.data)
	return nil
}

func (ms *MemStorage) GetCounter(name string) (int64, error) {
	ms.counter.mutex.Lock()
	defer ms.counter.mutex.Unlock()
	logger.Debugf("memstorage before Get counter: %#v", ms.counter.data)
	value, ok := ms.counter.data[name]
	if !ok {
		return 0, fmt.Errorf("counter %s does not exist", name)
	}
	return value, nil
}

func (ms *MemStorage) GetGauge(name string) (float64, error) {
	ms.gauge.mutex.Lock()
	defer ms.gauge.mutex.Unlock()
	logger.Debugf("memstorage before Get gauge: %#v", ms.gauge.data)
	value, ok := ms.gauge.data[name]
	if !ok {
		return 0, fmt.Errorf("gauge %s does not exist", name)
	}
	return value, nil
}
