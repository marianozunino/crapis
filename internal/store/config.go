package store

type StoreOption func(*Store)

func WithPassiveEvictionEnabled(enabled bool) StoreOption {
	return func(s *Store) {
		s.passiveEvictionEnabled = enabled
	}
}
