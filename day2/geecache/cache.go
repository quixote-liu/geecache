package geecache

import (
	"sync"

	"cache/day2/geecache/lru"
)

type cache struct {
	mu        sync.Mutex
	lru       *lru.Cache
	cacheBytes int64
}

func (c *cache) add(key string, value ByteViews) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.lru == nil {
		c.lru = lru.NewCache(c.cacheBytes, nil)
	}

	c.lru.Add(key, value)
}

func (c *cache) get(key string) (value ByteViews, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.lru == nil {
		c.lru = lru.NewCache(c.cacheBytes, nil)
	}

	if v, ok := c.lru.Get(key); ok {
		return v.(ByteViews), ok
	}

	return
}
