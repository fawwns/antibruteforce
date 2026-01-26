package bucket

import (
	"sync"
	"time"
)

type Store struct {
	mu      sync.Mutex
	Buckets map[string]*Bucket
}

func NewStore() *Store {
	return &Store{
		Buckets: make(map[string]*Bucket),
	}
}

// Get возвращает существующий бакет или создаёт новый.
func (s *Store) Get(key string, limit int) *Bucket {
	s.mu.Lock()
	defer s.mu.Unlock()

	if b, ok := s.Buckets[key]; ok {
		return b
	}

	b := &Bucket{
		Count:      0,
		Limit:      limit,
		LastRefill: time.Now(),
	}

	s.Buckets[key] = b
	return b
}

// Cleanup удаляет неактивные бакеты (старше N минут).
func (s *Store) Cleanup(ttl time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for key, val := range s.Buckets {
		if time.Since(val.LastRefill) > ttl {
			delete(s.Buckets, key)
		}
	}

}

func (s *Store) ResetAll() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Buckets = make(map[string]*Bucket)
}
