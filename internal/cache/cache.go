package cache

import (
	"sync"
	"time"

	"github.com/charlesaraya/pokedex-go/internal/pokedex"
)

type Cache struct {
	CachedEntries map[string]*CacheEntry
	Pokedex       *pokedex.Pokedex
	Mu            sync.RWMutex
}

type CacheEntry struct {
	CreatedAt time.Time
	Val       []byte
}

func NewCache(interval time.Duration) *Cache {
	var cache *Cache = &Cache{
		CachedEntries: make(map[string]*CacheEntry),
	}
	go cache.reapLoop(interval)
	return cache
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	// run every tick interval
	for range ticker.C {
		for k, v := range c.CachedEntries {
			// clean cached entries older than interval
			if time.Since(v.CreatedAt) > interval {
				delete(c.CachedEntries, k)
			}
		}
	}
}

func (c *Cache) Add(key string, val []byte) {
	c.Mu.Lock()
	c.CachedEntries[key] = &CacheEntry{
		CreatedAt: time.Now(),
		Val:       val,
	}
	defer c.Mu.Unlock()
}

func (c *Cache) Get(key string) (*CacheEntry, bool) {
	c.Mu.RLock()
	cachedEntry, ok := c.CachedEntries[key]
	defer c.Mu.RUnlock()
	return cachedEntry, ok
}
