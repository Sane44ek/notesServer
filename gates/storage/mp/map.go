package mp

import (
	"fmt"
	"notes/gates/storage"
	"reflect"
	"sync"
)

type Mp struct {
	
	body map[int64]interface{}
	mtx  sync.RWMutex
}

func NewMap() *Mp {
	return &Mp{body: nil}
}

func (m *Mp) Add(data interface{}) (key int64, e error) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	if m.body == nil {
		m.body = make(map[int64]interface{})
		m.body[0] = data
		return 0, nil
	}
	tp := reflect.TypeOf(m.body[0])
	if tp != reflect.TypeOf(data) {
		return -1, storage.ErrMismatchType
	}
	m.body[int64(len(m.body))] = data
	return int64(len(m.body)), nil
}

func (m *Mp) Print() {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	if m.body == nil {
		fmt.Println("no data")
		return
	}
	fmt.Println(m.body)
}

// Len возвращает длину списка
func (m *Mp) Len() (length int64) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	if m.body == nil {
		fmt.Println("no data")
		return
	}
	length = int64(len(m.body))
	return length
}

// RemoveByIndex удаляет элемент из списка по индексу
func (m *Mp) RemoveByIndex(id int64) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	if m.body == nil {
		fmt.Println("no data")
		return
	}
	delete(m.body, id)
}

// RemoveByValue удаляет элемент из списка по значению
func (m *Mp) RemoveByValue(value interface{}) bool {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	if m.body == nil {
		fmt.Println("no data")
		return false
	}
	for key := range m.body {
		if m.body[key] == value {
			delete(m.body, key)
			return true
		}
	}
	return false
}

func (m *Mp) RemoveAllByValue(value interface{}) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	if m.body == nil {
		fmt.Println("no data")
		return
	}
	for key := range m.body {
		if m.body[key] == value {
			delete(m.body, key)
		}
	}
}

func (m *Mp) GetByIndex(id int64) (value interface{}, ok bool) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	if m.body == nil {
		fmt.Println("no data")
		return 0, false
	}
	i, ok := m.body[id]
	return i, ok
}

func (m *Mp) GetByValue(value interface{}) (index int64, ok bool) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	if m.body == nil {
		fmt.Println("no data")
		return -1, false
	}
	for key := range m.body {
		if m.body[key] == value {
			return key, true
		}
	}
	return -1, false
}

func (m *Mp) GetAllByValue(value interface{}) (ids []int64, ok bool) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	if m.body == nil {
		fmt.Println("no data")
		return []int64{}, false
	}
	for key := range m.body {
		if m.body[key] == value {
			ids = append(ids, key)
		}
	}
	if len(ids) > 0 {
		return ids, true
	}
	return []int64{}, false
}

func (m *Mp) GetAll() (values []interface{}, ok bool) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	if m.body == nil {
		fmt.Println("no data")
		return []interface{}{}, false
	}
	for key := range m.body {
		values = append(values, m.body[key])
	}
	if len(values) > 0 {
		return values, true
	}
	return []interface{}{}, false
}

// Clear очищает список
func (m *Mp) Clear() {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.body = nil
}
