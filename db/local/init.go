package local

import "sync"

type Storage struct {
	hashMap map[string]struct{}
	mu      *sync.Mutex
}

func New() *Storage {
	return &Storage{
		hashMap: make(map[string]struct{}),
		mu:      &sync.Mutex{},
	}
}

func (s *Storage) AddHash(hash string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.hashMap[hash] = struct{}{}
}

func (s *Storage) IsHashExist(hash string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.hashMap[hash]
	if ok {
		return true
	}

	return false
}
