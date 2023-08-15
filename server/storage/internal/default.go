package internal

import (
	"fmt"
	"sync"

	"github.com/panupakm/miniredis/internal/payload"
)

type DefaultStorage struct {
	pairs map[string]payload.General
	mu    sync.RWMutex
}

func NewStorage() *DefaultStorage {
	return &DefaultStorage{
		pairs: make(map[string]payload.General),
	}
}

func NewStorageWithPair(pairs map[string]payload.General) *DefaultStorage {
	return &DefaultStorage{
		pairs: pairs,
	}
}

func (s *DefaultStorage) Set(key string, value payload.General) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.pairs[key] = value
	return nil
}

func (s *DefaultStorage) Get(key string) (payload.General, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	v, ok := s.pairs[key]
	if !ok {
		return payload.General{}, fmt.Errorf("key %s not found", key)
	}
	return v, nil
}
