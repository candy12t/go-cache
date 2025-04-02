package cache

import (
	"sync"
	"time"
)

type CacheItem struct {
	value  any
	expiry time.Time
}

type Cache struct {
	mu   sync.RWMutex
	data map[string]CacheItem
}

func NewCache() *Cache {
	return &Cache{
		data: make(map[string]CacheItem),
	}
}

func (c *Cache) Set(key string, value any, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var expiry time.Time
	if ttl > 0 {
		expiry = time.Now().Add(ttl)
	} else {
		expiry = time.Time{}
	}

	c.data[key] = CacheItem{
		value:  value,
		expiry: expiry,
	}
}

func (c *Cache) Get(key string) (any, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, ok := c.data[key]
	if !ok {
		return nil, false
	}

	if !item.expiry.IsZero() && item.expiry.Before(time.Now()) {
		delete(c.data, key)
		return nil, false
	}

	return item.value, ok
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.data, key)
}

func (c *Cache) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data = make(map[string]CacheItem)
}
