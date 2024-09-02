package resp

import (
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

type data struct {
	val string
	ttl *time.Time
}

type Store struct {
	setMap map[string]data
	mu     sync.RWMutex
}

func newStore() *Store {
	return &Store{
		setMap: map[string]data{},
		mu:     sync.RWMutex{},
	}
}

func (s *Store) StoreValue(key, val string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.setMap[key] = data{
		val: val,
	}
}

func (s *Store) ReadVal(key string) *string {
	s.mu.Lock()
	defer s.mu.Unlock()

	val, ok := s.setMap[key]
	if !ok {
		return nil
	}

	if val.ttl != nil {
		if val.ttl.Before(time.Now()) {
			log.Debug().Msgf("Passive Expiration of [%s]", key)
			delete(s.setMap, key)
			return nil
		}
	}
	return &val.val
}

func (s *Store) StoreValueWithTTL(key, val string, ttl int64) {
	now := time.Now()
	ttlDuration := time.Second * time.Duration(ttl)
	expireAt := now.Add(ttlDuration)

	s.mu.Lock()
	defer s.mu.Unlock()
	s.setMap[key] = data{
		val: val,
		ttl: &expireAt,
	}
}
