package cache

import (
	"time"

	memCache "github.com/patrickmn/go-cache"
)

type MemoryCache struct {
	cache *memCache.Cache
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		cache: memCache.New(memCache.NoExpiration, 5*time.Minute),
	}
}

func (c *MemoryCache) Get(key string) ([]byte, error) {
	v, ok := c.cache.Get(key)
	if !ok {
		return nil, nil
	}

	return v.([]byte), nil
}

func (c *MemoryCache) Set(key string, value []byte, ttl time.Duration) error {
	c.cache.Set(key, value, ttl)
	return nil
}

func (c *MemoryCache) Del(key string) error {
	c.cache.Delete(key)
	return nil
}

func (c *MemoryCache) Ping() error {
	return nil
}
