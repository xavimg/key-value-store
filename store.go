package main

import (
	"fmt"
	"sync"
)

type Storer[K comparable, V any] interface {
	Push(K, V) error
	Get(K) (V, error)
	Update(K, V) error
	Delete(K) error
}

type KVStore[K comparable, V any] struct {
	mu sync.RWMutex

	data map[K]V
}

func NewKVStore[K comparable, V any]() *KVStore[K, V] {
	return &KVStore[K, V]{
		data: make(map[K]V),
	}
}

func (s *KVStore[K, V]) Push(key K, value V) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = value

	return nil
}

func (s *KVStore[K, V]) Get(key K) (V, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, ok := s.data[key]
	if !ok {
		return value, fmt.Errorf("key-not-found")
	}

	return value, nil
}

func (s *KVStore[K, V]) Update(key K, value V) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if !s.has(key) {
		return nil
	}

	s.data[key] = value

	return nil
}

func (s *KVStore[K, V]) Delete(key K) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if !s.has(key) {
		return nil
	}

	delete(s.data, key)

	return nil
}

// NOTE: This is not concurrent safe, should be used with a Lock.
func (s *KVStore[K, V]) has(key K) bool {
	_, ok := s.data[key]
	if !ok {
		return false
	}

	return ok
}
