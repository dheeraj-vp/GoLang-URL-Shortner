package mock

import (
	"context"
	"fmt"
	"sync"
)

// ImprovedMockCache is an enhanced mock implementation of the Cache interface for testing
type ImprovedMockCache struct {
	data      map[string]string
	mu        sync.RWMutex
	getCount  int
	setCount  int
	delCount  int
	shouldFail bool
}

// NewImprovedMockCache creates a new improved mock cache instance
func NewImprovedMockCache() *ImprovedMockCache {
	return &ImprovedMockCache{
		data: make(map[string]string),
	}
}

// Set stores a key-value pair in the mock cache
func (m *ImprovedMockCache) Set(ctx context.Context, key string, val string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.shouldFail {
		return fmt.Errorf("mock cache: set operation failed")
	}

	m.setCount++
	m.data[key] = val
	return nil
}

// Get retrieves a value from the mock cache
func (m *ImprovedMockCache) Get(ctx context.Context, key string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.shouldFail {
		return "", fmt.Errorf("mock cache: get operation failed")
	}

	m.getCount++
	val, exists := m.data[key]
	if !exists {
		return "", nil // Return empty string for cache miss (like Redis)
	}
	return val, nil
}

// Delete removes a key from the mock cache
func (m *ImprovedMockCache) Delete(ctx context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.shouldFail {
		return fmt.Errorf("mock cache: delete operation failed")
	}

	m.delCount++
	delete(m.data, key)
	return nil
}

// SetFailureMode enables or disables failure simulation
func (m *ImprovedMockCache) SetFailureMode(fail bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.shouldFail = fail
}

// GetSetCount returns the number of Set operations performed
func (m *ImprovedMockCache) GetSetCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.setCount
}

// GetGetCount returns the number of Get operations performed
func (m *ImprovedMockCache) GetGetCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.getCount
}

// GetDeleteCount returns the number of Delete operations performed
func (m *ImprovedMockCache) GetDeleteCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.delCount
}

// Clear removes all data from the cache
func (m *ImprovedMockCache) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data = make(map[string]string)
	m.getCount = 0
	m.setCount = 0
	m.delCount = 0
}

// Size returns the number of items in the cache
func (m *ImprovedMockCache) Size() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.data)
}
