package cache

import (
	"sync"
	"time"
)

type CacheItem[T any] struct {
	value  *T
	expiry time.Time
}

type Cache[K comparable, T any] struct {
	mu   sync.RWMutex
	data map[K]CacheItem[T]
}

func NewCache[K comparable, T any]() *Cache[K, T] {
	return &Cache[K, T]{
		data: make(map[K]CacheItem[T]),
	}
}

func (c *Cache[K, T]) Set(key K, value *T, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var expiry time.Time
	if ttl > 0 {
		expiry = time.Now().Add(ttl)
	} else {
		expiry = time.Time{}
	}

	c.data[key] = CacheItem[T]{
		value:  value,
		expiry: expiry,
	}
}

func (c *Cache[K, T]) Get(key K) (*T, bool) {
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

func (c *Cache[K, T]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.data, key)
}

func (c *Cache[K, T]) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data = map[K]CacheItem[T]{}
}
