package store

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog/log"
)

type data struct {
	val string
	ttl time.Time
}

type Store struct {
	setMap  map[string]data
	ttlKeys map[string]struct{}
	mu      sync.RWMutex

	passiveEvictionEnabled atomic.Bool
	evictionIntervalMs     time.Duration
	evictionTimeoutMs      time.Duration
	stopChan               chan struct{}
}

func NewStore(opts ...StoreOption) *Store {
	s := &Store{
		setMap:             make(map[string]data),
		ttlKeys:            make(map[string]struct{}),
		mu:                 sync.RWMutex{},
		evictionIntervalMs: 250 * time.Millisecond,
		evictionTimeoutMs:  10 * time.Millisecond,
		stopChan:           make(chan struct{}),
	}
	s.passiveEvictionEnabled.Store(true)

	for _, opt := range opts {
		opt(s)
	}

	// evictionTimeoutMs must be at most half of evictionIntervalMs
	if s.evictionTimeoutMs > s.evictionIntervalMs/2 {
		log.Debug().Msg("evictionTimeoutMs must be at most half of evictionIntervalMs")
		s.evictionTimeoutMs = s.evictionIntervalMs / 2
		log.Debug().Msgf("evictionTimeoutMs set to half of evictionIntervalMs (%s)", s.evictionTimeoutMs)
	}

	if s.passiveEvictionEnabled.Load() {
		go s.startTTLExpirationThread()
	}

	return s
}

func (s *Store) startTTLExpirationThread() {
	log.Debug().Msgf("Starting Eviction Worker (interval=%s, timeout=%s)", s.evictionIntervalMs, s.evictionTimeoutMs)
	ticker := time.NewTicker(s.evictionIntervalMs)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// NOTE: the job can only take evictionTimeoutMs to avoid blocking other operations
			ctx, cancel := context.WithTimeout(context.Background(), s.evictionTimeoutMs)
			evictedKeys, err := s.deleteExpiredKeys(ctx)
			cancel()
			if err != nil {
				log.Debug().Msgf("Error deleting keys: %v", err)
				continue
			}
			if evictedKeys > 0 {
				log.Debug().Msgf("Deleted %d expired keys", evictedKeys)
			}
		case <-s.stopChan:
			log.Debug().Msgf("Stopping Eviction Worker")
			return
		}
	}
}

func (s *Store) Shutdown() {
	if s.passiveEvictionEnabled.Load() {
		close(s.stopChan)
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
				s.mu.Lock()
				delete(s.setMap, k)
				delete(s.ttlKeys, k)
				s.mu.Unlock()
				deletedKeys++
				log.Debug().Str("key", k).Msg("Active Expiration")
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
