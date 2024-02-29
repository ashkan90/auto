package utils

import (
	"log"
	"sync"
	"sync/atomic"
)

type SyncMap struct {
	count *atomic.Int64
	_map  *sync.Map
}

func NewSyncMap() *SyncMap {
	return &SyncMap{
		count: &atomic.Int64{},
		_map:  &sync.Map{},
	}
}

func (m *SyncMap) Len() int {
	if m.count == nil {
		return 0
	}

	ln := m.count.Load()

	return int(ln)
}

func (m *SyncMap) Get(key string) (any, bool) {
	if m._map == nil {
		return nil, false
	}

	return m._map.Load(key)
}

func (m *SyncMap) Add(key string, val any) {
	if m._map == nil {
		return
	}

	log.Println("[SyncMap] adding a value", key)

	m._map.Store(key, val)
	m.count.Add(1)
}

func (m *SyncMap) Delete(key string) {
	if m._map == nil {
		return
	}

	_, loaded := m._map.LoadAndDelete(key)
	if loaded {
		m.count.Add(-1)
	}
}

func (m *SyncMap) Range(f func(key, value any) bool) {
	if m._map == nil {
		return
	}

	m._map.Range(f)
}
