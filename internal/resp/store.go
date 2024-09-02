package resp

import "sync"

type Store struct {
	setMap map[string]string
	mu     sync.RWMutex
}

func newStore() *Store {
	return &Store{
		setMap: map[string]string{},
		mu:     sync.RWMutex{},
	}
}

func (s *Store) StoreValue(key, val string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.setMap[key] = val
}

func (s *Store) ReadVal(key string) *string {
	val, ok := s.setMap[key]
	if !ok {
		return nil
	}
	return &val

}
