package internal

import (
	"sync"
	"time"
)

type CacheEntry struct {
	CreatedAt time.Time
	Val       []byte
}

type Cache struct {
	Entries map[string]CacheEntry
	Mu      sync.Mutex
}

func (c *Cache) Add(key string, val []byte) {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	c.Entries[key] = CacheEntry{time.Now(), val}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	if entry, ok := c.Entries[key]; ok {
		return entry.Val, ok
	}
	return []byte{}, false
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		c.Mu.Lock()

		for key, entry := range c.Entries {
			age := time.Since(entry.CreatedAt)
			if age > interval {
				delete(c.Entries, key)
			}
		}

		c.Mu.Unlock()
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
