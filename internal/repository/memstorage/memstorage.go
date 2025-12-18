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
	if _, ok := ms.counter.data[name]; ok {
		ms.counter.mutex.Lock()
		ms.counter.data[name] += value
		ms.counter.mutex.Unlock()
		return nil
	}

	ms.counter.data[name] += value
	return nil
}

func (ms *MemStorage) PrintCounter() string {
	return fmt.Sprintf("%#v", storage.counter.data)
}

func (ms *MemStorage) SetGauge(name string, value float64) error {
	if _, ok := ms.gauge.data[name]; ok {
		ms.gauge.mutex.Lock()
		ms.gauge.data[name] = value
		ms.gauge.mutex.Unlock()
		return nil
	}

	ms.gauge.data[name] = value
	return nil
}

func (ms *MemStorage) PrintGauge() string {
	return fmt.Sprintf("%#v", storage.gauge.data)
}
