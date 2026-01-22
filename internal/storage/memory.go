package storage

import (
	"sync"
)

type MemoryStorage struct {
	mu   sync.RWMutex
	data map[string][]byte
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		data: make(map[string][]byte),
	}
}

func (s *MemoryStorage) Set(key string, value []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	result := make([]byte, len(value))
	copy(result, value)
	s.data[key] = result
	return nil
}

func (s *MemoryStorage) Get(key string) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, ok := s.data[key]
	if !ok {
		return nil, nil
	}
	result := make([]byte, len(value))
	copy(result, value)
	return result, nil
}

func (s *MemoryStorage) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
	return nil
}
