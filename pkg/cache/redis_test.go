package cache

import (
	"os"
	"testing"

	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
)

func newRedisTestClient() (*redis.Client, error) {
	var client *redis.Client
	if addr := os.Getenv("REDIS_ADDRESS"); addr != "" {
		client = redis.NewClient(&redis.Options{
			Addr: addr,
			DB:   2,
		})
	} else {
		client = redis.NewClient(&redis.Options{
			Addr: "127.0.0.1:6379",
			DB:   2,
		})
		if err := client.Ping().Err(); err != nil {
			return nil, nil
		}
	}
	if client != nil {
		client.FlushDB()
	}
	return client, nil
}

func TestRedisCache(t *testing.T) {
	client, err := newRedisTestClient()
	if err != nil {
		t.Fatal(err)
	} else if client == nil {
		t.Skip("no redis server available")
	}
	testCache(t, func() Cache {
		assert.NoError(t, client.FlushDB().Err())
		return &RedisCache{
			Client: client,
		}
	})
}
