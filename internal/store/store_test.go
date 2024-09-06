package store

import (
	"context"
	"testing"
	"time"
)

func TestStore(t *testing.T) {
	t.Run("StoreValue and ReadVal", func(t *testing.T) {
		tests := []struct {
			name     string
			key      string
			value    *string
			expected *string
		}{
			{"Store and read value", "key1", ptr("value1"), ptr("value1")},
			{"Read non-existent key", "key3", nil, nil},
		}

		for _, tt := range tests {
			store := NewStore()
			t.Run(tt.name, func(t *testing.T) {
				if tt.value != nil {
					store.StoreValue(tt.key, *tt.value)
				}
				result := store.ReadValue(tt.key)
				if !equalPtrString(result, tt.expected) {
					t.Errorf("Expected %v, got %v", tt.expected, result)
				}
			})
		}
	})

	t.Run("StoreValueWithTTL", func(t *testing.T) {
		tests := []struct {
			name        string
			key         string
			value       string
			ttl         int64
			sleepBefore time.Duration
			expected    *string
		}{
			{"Store with TTL and read before expiration", "key1", "value1", 5, 1 * time.Millisecond, ptr("value1")},
			{"Store with TTL and read after expiration", "key2", "value2", 1, 1_100 * time.Millisecond, nil},
		}

		store := NewStore()
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				store.StoreValueWithTTL(tt.key, tt.value, tt.ttl)
				time.Sleep(tt.sleepBefore)
				result := store.ReadValue(tt.key)
				if !equalPtrString(result, tt.expected) {
					t.Errorf("Expected %v, got %v", tt.expected, result)
				}
			})
		}
	})

	t.Run("DeleteKey", func(t *testing.T) {
		tests := []struct {
			name            string
			initialKeys     map[string]string
			keysToDelete    []string
			expectedCount   int
			remainingKeys   []string
			nonExistentKeys []string
		}{
			{
				name:            "Delete existing keys",
				initialKeys:     map[string]string{"key1": "value1", "key2": "value2", "key3": "value3"},
				keysToDelete:    []string{"key1", "key3"},
				expectedCount:   2,
				remainingKeys:   []string{"key2"},
				nonExistentKeys: []string{"key1", "key3"},
			},
			{
				name:            "Delete non-existent keys",
				initialKeys:     map[string]string{"key1": "value1"},
				keysToDelete:    []string{"key2", "key3"},
				expectedCount:   0,
				remainingKeys:   []string{"key1"},
				nonExistentKeys: []string{"key2", "key3"},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				store := NewStore()
				for k, v := range tt.initialKeys {
					store.StoreValue(k, v)
				}

				count := store.DeleteKeys(tt.keysToDelete...)
				if count != tt.expectedCount {
					t.Errorf("Expected %d deletions, got %d", tt.expectedCount, count)
				}

				for _, k := range tt.remainingKeys {
					if store.ReadValue(k) == nil {
						t.Errorf("Key %s should still exist", k)
					}
				}

				for _, k := range tt.nonExistentKeys {
					if store.ReadValue(k) != nil {
						t.Errorf("Key %s should not exist", k)
					}
				}
			})
		}
	})

	t.Run("TTL Expiration", func(t *testing.T) {
		store := NewStore()
		store.StoreValueWithTTL("key1", "value1", 1) // 1 second TTL

		// Wait for TTL to expire
		time.Sleep(1_250 * time.Millisecond)

		// Check if the key has been removed
		if store.ReadValue("key1") != nil {
			t.Errorf("Key 'key1' should have been removed due to TTL expiration")
		}
	})

	t.Run("deleteExpiredKeys", func(t *testing.T) {
		store := NewStore(WithPassiveEviction(false))

		store.StoreValueWithTTL("key1", "value1", 1) // 1 second TTL
		store.StoreValueWithTTL("key2", "value2", 2) // 3 seconds TTL
		store.StoreValue("key3", "value3")           // No TTL

		// Wait for the first key to expire
		time.Sleep(1_250 * time.Millisecond)

		ctx := context.Background()
		deletedKeys, err := store.deleteExpiredKeys(ctx)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if deletedKeys != 1 {
			t.Errorf("Expected 1 deleted key, got %d", deletedKeys)
		}

		// Check remaining keys
		if store.ReadValue("key1") != nil {
			t.Errorf("Key 'key1' should have been deleted")
		}
		if store.ReadValue("key2") == nil {
			t.Errorf("Key 'key2' should still exist")
		}
		if store.ReadValue("key3") == nil {
			t.Errorf("Key 'key3' should still exist")
		}
	})
}

// Helper function to create a pointer to a string
func ptr(s string) *string {
	return &s
}

// Helper function to compare two string pointers
func equalPtrString(a, b *string) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}
