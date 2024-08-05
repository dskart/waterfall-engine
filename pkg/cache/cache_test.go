package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func testCache(t *testing.T, factory func() Cache) {
	t.Run("GetAndSet", func(t *testing.T) {
		c := factory()

		assert.NoError(t, c.Set("foo", []byte("bar"), time.Minute))
		v, err := c.Get("foo")
		assert.NoError(t, err)
		assert.Equal(t, []byte("bar"), v)

		didExpire := false
		assert.NoError(t, c.Set("foo", []byte("bar"), 10*time.Millisecond))
		for i := 0; i < 10; i++ {
			v, err := c.Get("foo")
			assert.NoError(t, err)
			if v == nil {
				didExpire = true
				break
			}
			time.Sleep(5 * time.Millisecond)
		}
		assert.True(t, didExpire)
	})

	t.Run("Del", func(t *testing.T) {
		b := factory()

		err := b.Del("foo")
		assert.NoError(t, err)

		assert.NoError(t, b.Set("foo", []byte("bar"), time.Minute))
		v, err := b.Get("foo")
		assert.NotNil(t, v)
		assert.NoError(t, err)

		err = b.Del("foo")
		assert.NoError(t, err)
		v, err = b.Get("foo")
		assert.Nil(t, v)
		assert.NoError(t, err)
	})
}
