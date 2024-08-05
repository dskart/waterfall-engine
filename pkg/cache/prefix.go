package cache

import (
	"encoding/binary"
	"time"
)

// NameSpacedCache returns a cache that prefixes all keys with a namespace.
// This is useful for when you want to use the same cache and avoid key conflicts.
func NameSpacedCache(c Cache, namespace string) Cache {
	buf := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(buf, uint64(len(namespace)))
	return &PrefixCache{
		Prefix: "@" + string(buf[:n]) + namespace,
		Cache:  c,
	}
}

type PrefixCache struct {
	Prefix string
	Cache  Cache
}

func (c *PrefixCache) Get(key string) ([]byte, error) {
	return c.Cache.Get(c.Prefix + key)
}

func (c *PrefixCache) Set(key string, value []byte, ttl time.Duration) error {
	return c.Cache.Set(c.Prefix+key, value, ttl)
}

func (c *PrefixCache) Del(key string) error {
	return c.Cache.Del(c.Prefix + key)
}

func (c *PrefixCache) Ping() error {
	return c.Cache.Ping()
}
