package collections

import (
	"sync"
)

// BeeMap is a map with lock
type Map struct {
	sync.RWMutex
	bm map[interface{}]interface{}
}

// NewBeeMap return new safemap
func NewMap() *Map {
	return &Map{
		bm: make(map[interface{}]interface{}),
	}
}

func (m *Map) Add(k interface{}, v interface{}) {
	m.Lock()
	defer m.Unlock()
	if val, ok := m.bm[k]; !ok {
		m.bm[k] = v
	} else if val != v {
		m.bm[k] = v
	}
}

// Get from maps return the k's value
func (m *Map) Get(k interface{}) interface{} {
	m.RLock()
	defer m.RUnlock()
	if val, ok := m.bm[k]; ok {
		return val
	}
	return nil
}

// Set Maps the given key and value. Returns false
// if the key is already in the map and changes nothing.
func (m *Map) Set(k interface{}, v interface{}) bool {
	m.Lock()
	defer m.Unlock()
	if val, ok := m.bm[k]; !ok {
		m.bm[k] = v
	} else if val != v {
		m.bm[k] = v
	} else {
		return false
	}
	return true
}

// Has Returns true if k is exist in the map.
func (m *Map) Has(k interface{}) bool {
	m.RLock()
	defer m.RUnlock()
	_, ok := m.bm[k]
	return ok
}

// Delete the given key and value.
func (m *Map) Remove(k interface{}) {
	m.Lock()
	defer m.Unlock()
	delete(m.bm, k)
}

// Items returns all items in safemap.
func (m *Map) Items() map[interface{}]interface{} {
	m.RLock()
	defer m.RUnlock()
	r := make(map[interface{}]interface{})
	for k, v := range m.bm {
		r[k] = v
	}
	return r
}

// Count returns the number of items within the map.
func (m *Map) Len() int {
	m.RLock()
	defer m.RUnlock()
	return len(m.bm)
}

func (m *Map) IsEmpty() bool {
	m.RLock()
	defer m.RUnlock()
	if len(m.bm) == 0 {
		return true
	}
	return false
}
