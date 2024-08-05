package cache

import (
	"testing"
)

func TestMemoryCache(t *testing.T) {
	testCache(t, func() Cache {
		return NewMemoryCache()
	})
}
