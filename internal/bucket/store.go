package bucket

import (
	"sync"
	"time"
)

type Store struct {
	mu      sync.Mutex
	buckets map[string]*Bucket
}

func NewStore() *Store {
	return &Store{
		buckets: make(map[string]*Bucket),
	}
}

// Get возвращает существующий бакет или создаёт новый.
func (s *Store) Get(key string, limit int) *Bucket {
	s.mu.Lock()
	defer s.mu.Unlock()

	if b, ok := s.buckets[key]; ok {
		return b
	}

	b := &Bucket{
		Count:      0,
		Limit:      limit,
		LastRefill: time.Now(),
	}

	s.buckets[key] = b
	return b
}

// Cleanup удаляет неактивные бакеты (старше N минут).
func (s *Store) Cleanup(ttl time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for key, val := range s.buckets {
		if time.Since(val.LastRefill) > ttl {
			delete(s.buckets, key)
		}
	}

}
