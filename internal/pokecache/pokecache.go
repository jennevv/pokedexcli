package pokecache

import (
	"sync"
	"time"
)

type CacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	entries map[string]CacheEntry
	mu      sync.Mutex
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[key] = CacheEntry{time.Now(), val}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if entry, ok := c.entries[key]; ok {
		return entry.val, ok
	}
	return []byte{}, false
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		c.mu.Lock()

		for key, entry := range c.entries {
			age := time.Since(entry.createdAt)
			if age > interval {
				delete(c.entries, key)
			}
		}

		c.mu.Unlock()
	}
}

func NewCache(interval time.Duration) *Cache {
	cache := Cache{
		map[string]CacheEntry{},
		sync.Mutex{},
	}

	go cache.reapLoop(interval)

	return &cache
}
