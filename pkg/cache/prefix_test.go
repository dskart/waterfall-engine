package cache

import "testing"

func TestPrefixCache(t *testing.T) {
	testCache(t, func() Cache {
		return &PrefixCache{
			Prefix: "foo",
			Cache:  NewMemoryCache(),
		}
	})
}
