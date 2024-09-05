package store

import "time"

type StoreOption func(*Store)

func WithEvictionInterval(interval time.Duration) StoreOption {
	return func(s *Store) {
		s.evictionIntervalMs = interval
	}
}

func WithEvictionTimeout(timeout time.Duration) StoreOption {
	return func(s *Store) {
		s.evictionTimeoutMs = timeout
	}
}

func WithPassiveEviction(enabled bool) StoreOption {
	return func(s *Store) {
		s.passiveEvictionEnabled.Store(enabled)
	}
}
