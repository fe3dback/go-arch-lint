package safemap

import "sync"

type Safemap[K comparable, T any] struct {
	data map[K]T
	mux  sync.RWMutex
}

func New[K comparable, T any]() *Safemap[K, T] {
	return &Safemap[K, T]{
		data: map[K]T{},
	}
}

func (m *Safemap[K, T]) Set(key K, value T) {
	m.mux.Lock()
	m.data[key] = value
	m.mux.Unlock()
}

func (m *Safemap[K, T]) Get(key K) (T, bool) {
	m.mux.RLock()
	value, exist := m.data[key]
	m.mux.RUnlock()

	return value, exist
}

func (m *Safemap[K, T]) Iterate(itt func(K, T)) {
	m.mux.RLock()
	for k, t := range m.data {
		itt(k, t)
	}
	m.mux.RUnlock()
}
