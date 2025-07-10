package cache

import (
	"sync"
	"time"

	"github.com/gouthamve/hardcover-book-embed/internal/hardcover"
)

type CacheItem struct {
	Data      *hardcover.CurrentlyReadingResponse
	ExpiresAt time.Time
}

type MemoryCache struct {
	mu    sync.RWMutex
	items map[string]*CacheItem
	ttl   time.Duration
}

func NewMemoryCache(ttl time.Duration) *MemoryCache {
	cache := &MemoryCache{
		items: make(map[string]*CacheItem),
		ttl:   ttl,
	}

	go cache.cleanup()
	return cache
}

func (c *MemoryCache) Get(key string) (*hardcover.CurrentlyReadingResponse, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return nil, false
	}

	if time.Now().After(item.ExpiresAt) {
		return nil, false
	}

	return item.Data, true
}

func (c *MemoryCache) Set(key string, data *hardcover.CurrentlyReadingResponse) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = &CacheItem{
		Data:      data,
		ExpiresAt: time.Now().Add(c.ttl),
	}
}

func (c *MemoryCache) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for key, item := range c.items {
			if now.After(item.ExpiresAt) {
				delete(c.items, key)
			}
		}
		c.mu.Unlock()
	}
}
