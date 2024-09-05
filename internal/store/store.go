package store

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

type data struct {
	val string
	ttl time.Time
}

type Store struct {
	setMap  map[string]data
	ttlKeys map[string]struct{} // NOTE: Saving 1 byte 🎉
	mu      sync.RWMutex

	passiveEvictionEnabled bool
	evectionIntervalMs     time.Duration
	evictionTimeoutMs      time.Duration
}

func NewStore(opts ...StoreOption) *Store {
	s := &Store{
		setMap:  make(map[string]data),
		ttlKeys: make(map[string]struct{}),
		mu:      sync.RWMutex{},

		passiveEvictionEnabled: true,
		evectionIntervalMs:     250 * time.Millisecond,
		evictionTimeoutMs:      10 * time.Millisecond,
	}

	for _, opt := range opts {
		opt(s)
	}

	// evictionTimeoutMs must be at least half of evectionIntervalMs
	if s.evictionTimeoutMs < s.evectionIntervalMs/2 {
		s.evictionTimeoutMs = s.evectionIntervalMs / 2
	}

	if s.passiveEvictionEnabled {
		go s.startTTLExpirationThread()
	}

	return s
}

func (s *Store) startTTLExpirationThread() {
	log.Debug().Msg("Starting TTL Expiration Thread")
	// NOTE: run 4 times a second
	ticker := time.NewTicker(s.evectionIntervalMs)
	defer ticker.Stop()

	for range ticker.C {
		// NOTE: the job can only take 10ms to avoid blocking other operations
		ctx, cancel := context.WithTimeout(context.Background(), s.evictionTimeoutMs)
		evictedKeys, err := s.deleteExpiredKeys(ctx)
		cancel()

		if err != nil {
			log.Debug().Err(err).Msg("Error deleting keys")
			continue
		}
		if evictedKeys > 0 {
			log.Debug().Int("count", evictedKeys).Msg("Deleted expired keys")
		}
	}
}

func (s *Store) deleteExpiredKeys(ctx context.Context) (int, error) {
	now := time.Now()
	deletedKeys := 0

	for k := range s.ttlKeys {
		select {
		case <-ctx.Done():
			return deletedKeys, ctx.Err()
		default:
			if v, exists := s.setMap[k]; exists && !v.ttl.IsZero() && v.ttl.Before(now) {
				log.Debug().Str("key", k).Msg("Active Expiration")
				s.mu.Lock()
				delete(s.setMap, k)
				delete(s.ttlKeys, k)
				s.mu.Unlock()
				deletedKeys++
			}
		}
	}
	return deletedKeys, nil
}

func (s *Store) StoreValue(key, val string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.setMap[key] = data{val: val}
	delete(s.ttlKeys, key) // Remove from TTL registry if it was there
}

func (s *Store) ReadVal(key string) *string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if val, ok := s.setMap[key]; ok {
		if !val.ttl.IsZero() && val.ttl.Before(time.Now()) {
			log.Debug().Str("key", key).Msg("Passive Expiration")
			return nil
		}
		return &val.val
	}
	return nil
}

func (s *Store) StoreValueWithTTL(key, val string, ttl int64) {
	expireAt := time.Now().Add(time.Duration(ttl) * time.Second)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.setMap[key] = data{val: val, ttl: expireAt}
	s.ttlKeys[key] = struct{}{} // Add to TTL registry
}

func (s *Store) DeleteKey(keys ...string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	deletedKeys := 0
	for _, k := range keys {
		if _, ok := s.setMap[k]; ok {
			delete(s.setMap, k)
			delete(s.ttlKeys, k)
			deletedKeys++
		}
	}
	return deletedKeys
}

func (s *Store) Expire(key string, ttl int64) int {
	// find key and set ttl
	s.mu.Lock()
	defer s.mu.Unlock()

	if val, ok := s.setMap[key]; ok {
		expireAt := time.Now().Add(time.Duration(ttl) * time.Second)
		val.ttl = expireAt

		s.setMap[key] = val
		s.ttlKeys[key] = struct{}{} // Add to TTL registry

		return 1
	}

	return 0
}
