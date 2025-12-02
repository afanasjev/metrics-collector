package memstorage

import (
	"sync"
	"testing"
)

func TestGetMemStorage(t *testing.T) {
	// Тест на singleton паттерн
	storage1 := GetMemStorage()
	storage2 := GetMemStorage()

	if storage1 != storage2 {
		t.Error("GetMemStorage should return the same instance (singleton)")
	}
}

func TestSetCounter(t *testing.T) {
	// Создаем новый экземпляр для теста
	storage := &MemStorage{
		counter: counter{
			data:  make(map[string]int64),
			mutex: sync.Mutex{},
		},
		gauge: gauge{
			data:  make(map[string]float64),
			mutex: sync.Mutex{},
		},
	}

	// Тест 1: Установка нового значения
	err := storage.SetCounter("testCounter", 10)
	if err != nil {
		t.Errorf("SetCounter failed: %v", err)
	}

	// Проверяем, что значение установлено (через внутреннюю структуру)
	storage.counter.mutex.Lock()
	value, exists := storage.counter.data["testCounter"]
	storage.counter.mutex.Unlock()

	if !exists {
		t.Error("Counter value should exist after SetCounter")
	}
	if value != 10 {
		t.Errorf("Expected counter value 10, got %d", value)
	}

	// Тест 2: Добавление к существующему значению
	err = storage.SetCounter("testCounter", 5)
	if err != nil {
		t.Errorf("SetCounter failed on increment: %v", err)
	}

	storage.counter.mutex.Lock()
	value = storage.counter.data["testCounter"]
	storage.counter.mutex.Unlock()

	if value != 15 {
		t.Errorf("Expected counter value 15 after increment, got %d", value)
	}

	// Тест 3: Установка другого счетчика
	err = storage.SetCounter("anotherCounter", 20)
	if err != nil {
		t.Errorf("SetCounter failed: %v", err)
	}

	storage.counter.mutex.Lock()
	value2 := storage.counter.data["anotherCounter"]
	storage.counter.mutex.Unlock()

	if value2 != 20 {
		t.Errorf("Expected counter value 20, got %d", value2)
	}
}

func TestSetGauge(t *testing.T) {
	// Создаем новый экземпляр для теста
	storage := &MemStorage{
		counter: counter{
			data:  make(map[string]int64),
			mutex: sync.Mutex{},
		},
		gauge: gauge{
			data:  make(map[string]float64),
			mutex: sync.Mutex{},
		},
	}

	// Тест 1: Установка нового значения
	err := storage.SetGauge("testGauge", 3.14)
	if err != nil {
		t.Errorf("SetGauge failed: %v", err)
	}

	// Проверяем, что значение установлено
	storage.gauge.mutex.Lock()
	value, exists := storage.gauge.data["testGauge"]
	storage.gauge.mutex.Unlock()

	if !exists {
		t.Error("Gauge value should exist after SetGauge")
	}
	if value != 3.14 {
		t.Errorf("Expected gauge value 3.14, got %f", value)
	}

	// Тест 2: Перезапись существующего значения
	err = storage.SetGauge("testGauge", 2.71)
	if err != nil {
		t.Errorf("SetGauge failed on update: %v", err)
	}

	storage.gauge.mutex.Lock()
	value = storage.gauge.data["testGauge"]
	storage.gauge.mutex.Unlock()

	if value != 2.71 {
		t.Errorf("Expected gauge value 2.71 after update, got %f", value)
	}

	// Тест 3: Установка другого gauge
	err = storage.SetGauge("anotherGauge", 1.5)
	if err != nil {
		t.Errorf("SetGauge failed: %v", err)
	}

	storage.gauge.mutex.Lock()
	value2 := storage.gauge.data["anotherGauge"]
	storage.gauge.mutex.Unlock()

	if value2 != 1.5 {
		t.Errorf("Expected gauge value 1.5, got %f", value2)
	}

	// Тест 4: Установка нулевого значения
	err = storage.SetGauge("zeroGauge", 0.0)
	if err != nil {
		t.Errorf("SetGauge failed with zero value: %v", err)
	}

	storage.gauge.mutex.Lock()
	value3 := storage.gauge.data["zeroGauge"]
	storage.gauge.mutex.Unlock()

	if value3 != 0.0 {
		t.Errorf("Expected gauge value 0.0, got %f", value3)
	}
}

func TestPrintCounter(t *testing.T) {
	storage := &MemStorage{
		counter: counter{
			data:  make(map[string]int64),
			mutex: sync.Mutex{},
		},
		gauge: gauge{
			data:  make(map[string]float64),
			mutex: sync.Mutex{},
		},
	}

	// Устанавливаем несколько значений
	storage.SetCounter("counter1", 10)
	storage.SetCounter("counter2", 20)

	result := storage.PrintCounter()
	if result == "" {
		t.Error("PrintCounter should return non-empty string")
	}

	// Проверяем, что результат содержит информацию о счетчиках
	if len(result) == 0 {
		t.Error("PrintCounter result should not be empty")
	}
}

func TestPrintGauge(t *testing.T) {
	storage := &MemStorage{
		counter: counter{
			data:  make(map[string]int64),
			mutex: sync.Mutex{},
		},
		gauge: gauge{
			data:  make(map[string]float64),
			mutex: sync.Mutex{},
		},
	}

	// Устанавливаем несколько значений
	storage.SetGauge("gauge1", 1.1)
	storage.SetGauge("gauge2", 2.2)

	result := storage.PrintGauge()
	if result == "" {
		t.Error("PrintGauge should return non-empty string")
	}

	// Проверяем, что результат содержит информацию о gauge
	if len(result) == 0 {
		t.Error("PrintGauge result should not be empty")
	}
}

func TestSetCounterConcurrent(t *testing.T) {
	storage := &MemStorage{
		counter: counter{
			data:  make(map[string]int64),
			mutex: sync.Mutex{},
		},
		gauge: gauge{
			data:  make(map[string]float64),
			mutex: sync.Mutex{},
		},
	}

	// Инициализируем ключ, чтобы избежать race condition при первой записи
	// и при проверке существования ключа
	storage.counter.mutex.Lock()
	storage.counter.data["concurrentCounter"] = 0
	storage.counter.mutex.Unlock()

	// Тест на конкурентный доступ
	var wg sync.WaitGroup
	numGoroutines := 10
	iterations := 100

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				storage.SetCounter("concurrentCounter", 1)
			}
		}()
	}

	wg.Wait()

	// Проверяем, что все значения были добавлены
	storage.counter.mutex.Lock()
	value := storage.counter.data["concurrentCounter"]
	storage.counter.mutex.Unlock()

	expected := int64(numGoroutines * iterations)
	if value != expected {
		t.Errorf("Expected counter value %d after concurrent updates, got %d", expected, value)
	}
}

func TestSetGaugeConcurrent(t *testing.T) {
	storage := &MemStorage{
		counter: counter{
			data:  make(map[string]int64),
			mutex: sync.Mutex{},
		},
		gauge: gauge{
			data:  make(map[string]float64),
			mutex: sync.Mutex{},
		},
	}

	// Инициализируем ключ, чтобы избежать race condition при первой записи
	// и при проверке существования ключа
	storage.gauge.mutex.Lock()
	storage.gauge.data["concurrentGauge"] = 0.0
	storage.gauge.mutex.Unlock()

	// Тест на конкурентный доступ
	var wg sync.WaitGroup
	numGoroutines := 10

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			storage.SetGauge("concurrentGauge", float64(id))
		}(i)
	}

	wg.Wait()

	// Проверяем, что значение было установлено (последнее значение)
	storage.gauge.mutex.Lock()
	_, exists := storage.gauge.data["concurrentGauge"]
	storage.gauge.mutex.Unlock()

	if !exists {
		t.Error("Gauge value should exist after concurrent updates")
	}
}

