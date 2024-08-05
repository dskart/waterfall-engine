package cache

import (
	"fmt"
	"time"
)

type Cache interface {
	// Get returns the value at key or nil if it does not exist.
	Get(key string) ([]byte, error)
	// Set sets the value at key and returns an error if it failed.
	Set(key string, value []byte, ttl time.Duration) error
	// Del deletes the value at key
	Del(key string) error
	// Ping checks the connection to the cache and returns an error if it failed.
	Ping() error
}

func New(cfg Config) (Cache, error) {
	if cfg.InMemory {
		return NewMemoryCache(), nil
	}

	if cfg.RedisAddress != "" {
		return NewRedisCache(cfg.RedisAddress), nil

	}

	return nil, fmt.Errorf("invalid cache configuration")
}
