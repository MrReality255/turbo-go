package utils

import "sync"

type SafeMap[K comparable, V any] struct {
	mx   sync.Mutex
	data map[K]V
}

func NewSafeMap[K comparable, V any]() *SafeMap[K, V] {
	return &SafeMap[K, V]{
		data: make(map[K]V),
	}
}

func (m *SafeMap[K, V]) Set(key K, value V) {
	m.mx.Lock()
	defer m.mx.Unlock()
	m.data[key] = value
}

func (m *SafeMap[K, V]) Get(key K) V {
	m.mx.Lock()
	defer m.mx.Unlock()
	return m.data[key]
}

func (m *SafeMap[K, V]) Delete(k K) {
	m.mx.Lock()
	defer m.mx.Unlock()
	delete(m.data, k)
}

func (m *SafeMap[K, V]) Prepare(id K, f func() V) {
	m.mx.Lock()
	defer m.mx.Unlock()

	if _, ok := m.data[id]; !ok {
		m.data[id] = f()
	}
}

func (m *SafeMap[K, V]) GetValues(lessFct func(v1 V, v2 V) bool) []V {
	m.mx.Lock()
	defer m.mx.Unlock()
	return MapValues(m.data, lessFct)
}

func (m *SafeMap[K, V]) Clone() map[K]V {
	m.mx.Lock()
	defer m.mx.Unlock()
	return MapClone(m.data)
}

func (m *SafeMap[K, V]) Contains(id K) bool {
	m.mx.Lock()
	defer m.mx.Unlock()

	_, ok := m.data[id]
	return ok
}
