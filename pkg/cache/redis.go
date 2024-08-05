package cache

import (
	"time"

	"github.com/go-redis/redis"
)

type RedisCache struct {
	Client *redis.Client
}

func NewRedisCache(addr string) *RedisCache {
	return &RedisCache{
		Client: redis.NewClient(&redis.Options{
			Addr: addr,
		}),
	}
}

func (c *RedisCache) Get(key string) ([]byte, error) {
	v, err := c.Client.Get(key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return []byte(v), err
}

func (c *RedisCache) Set(key string, value []byte, ttl time.Duration) error {
	return c.Client.Set(key, value, ttl).Err()
}

func (c *RedisCache) Del(key string) error {
	result := c.Client.Del(key)
	return result.Err()
}

func (c *RedisCache) Ping() error {
	return c.Client.Ping().Err()
}
