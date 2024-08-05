package cache

import (
	"math"
	"math/rand"
	"time"

	"github.com/vmihailenco/msgpack/v5"
)

type earlyExpirationValue struct {
	Delta  time.Duration `msgpack:"d"`
	Expiry time.Time     `msgpack:"e"`
	Value  []byte        `msgpack:"v"`
}

// XFetch implements probabilistic early expiration. The "X" in the name is for eXponential
// function.
//
// The lock ensures that only one caller recomputes the value at a time. If the value is already
// being recomputed by one goroutine, other goroutines will either wait for that recomputation to
// finish or return an existing value if one is available.
//
// This will cache the return value of recompute, even if it is nil. To prevent writing to the
// cache, recompute should return a 0 ttl.
//
// See http://cseweb.ucsd.edu/~avattani/papers/cache_stampede.pdf
func XFetch(c Cache, lock *Lock, key string, beta float64, recompute func() (value []byte, ttl time.Duration, err error)) ([]byte, error) {
	buf, err := c.Get(key)
	if err != nil {
		return nil, err
	}

	var existing *earlyExpirationValue
	if buf != nil {
		if err := msgpack.Unmarshal(buf, &existing); err != nil {
			return nil, err
		}
	}

	now := time.Now()
	if existing == nil || !now.Add(time.Duration(-float64(existing.Delta)*beta*math.Log(rand.Float64()))).Before(existing.Expiry) {
		isOwner := false
		l := lock.Acquire(key, func() (interface{}, error) {
			isOwner = true
			start := now
			value, ttl, err := recompute()
			if err != nil {
				return nil, err
			}
			now := time.Now()
			created := &earlyExpirationValue{
				Delta:  now.Sub(start),
				Expiry: now.Add(ttl),
				Value:  value,
			}
			if ttl > 0 {
				if buf, err := msgpack.Marshal(created); err != nil {
					return nil, err
				} else if err := c.Set(key, buf, ttl); err != nil {
					return nil, err
				}
			}
			return created, nil
		})
		if existing == nil || isOwner {
			created, err := l.Wait()
			if err != nil {
				return nil, err
			}
			existing = created.(*earlyExpirationValue)
		}
	}

	return existing.Value, nil
}
