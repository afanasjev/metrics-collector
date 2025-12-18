package memstorage

import (
	"fmt"
	"sync"
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

func GetMemStorage() *MemStorage {
	once.Do(func() {
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
	return nil
}

func (ms *MemStorage) PrintCounter() string {
	return fmt.Sprintf("%#v", storage.counter.data)
}

func (ms *MemStorage) SetGauge(name string, value float64) error {
	ms.gauge.mutex.Lock()
	defer ms.gauge.mutex.Unlock()
	ms.gauge.data[name] = value
	return nil
}

func (ms *MemStorage) PrintGauge() string {
	return fmt.Sprintf("%#v", storage.gauge.data)
}

func (ms *MemStorage) GetCounter(name string) (int64, error) {
	ms.counter.mutex.Lock()
	defer ms.counter.mutex.Unlock()
	value, ok := ms.counter.data[name]
	if !ok {
		return 0, fmt.Errorf("counter %s does not exist", name)
	}
	return value, nil
}

func (ms *MemStorage) GetGauge(name string) (float64, error) {
	ms.gauge.mutex.Lock()
	defer ms.gauge.mutex.Unlock()
	value, ok := ms.gauge.data[name]
	if !ok {
		return 0, fmt.Errorf("gauge %s does not exist", name)
	}
	return value, nil
}
