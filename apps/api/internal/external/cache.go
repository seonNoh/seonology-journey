package external

import (
	"context"
	"sync"
	"time"
)

// MemoryCache is a simple in-memory cache for development and testing.
type MemoryCache struct {
	mu    sync.RWMutex
	store map[string]cacheEntry
}

type cacheEntry struct {
	data      []byte
	expiresAt time.Time
}

// NewMemoryCache creates a memory cache.
func NewMemoryCache() *MemoryCache {
	return &MemoryCache{store: make(map[string]cacheEntry)}
}

// Get retrieves a cached value. Returns nil if not found or expired.
func (c *MemoryCache) Get(_ context.Context, key string) ([]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	e, ok := c.store[key]
	if !ok || time.Now().After(e.expiresAt) {
		return nil, nil
	}
	return e.data, nil
}

// Set stores a value with TTL.
func (c *MemoryCache) Set(_ context.Context, key string, data []byte, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[key] = cacheEntry{data: data, expiresAt: time.Now().Add(ttl)}
	return nil
}
