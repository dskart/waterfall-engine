package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestXFetch(t *testing.T) {
	c := NewMemoryCache()
	var l Lock

	ttl := 10 * time.Second

	v, err := XFetch(c, &l, "foo", 1.0, func() ([]byte, time.Duration, error) {
		time.Sleep(200 * time.Millisecond)
		return []byte("foo"), ttl, nil
	})
	require.NoError(t, err)
	assert.Equal(t, "foo", string(v))

	v, err = XFetch(c, &l, "foo", 1.0, func() ([]byte, time.Duration, error) {
		return []byte("foo should not be recomputed yet"), ttl, nil
	})
	require.NoError(t, err)
	assert.Equal(t, "foo", string(v))

	for {
		v, err := XFetch(c, &l, "foo", 5.0, func() ([]byte, time.Duration, error) {
			return []byte("recomputed"), ttl, nil
		})
		require.NoError(t, err)
		if string(v) == "recomputed" {
			break
		}
	}

	t.Run("Nil", func(t *testing.T) {
		v, err := XFetch(c, &l, "nil", 1.0, func() ([]byte, time.Duration, error) {
			return nil, ttl, nil
		})
		require.NoError(t, err)
		assert.Nil(t, v)

		v, err = XFetch(c, &l, "nil", 1.0, func() ([]byte, time.Duration, error) {
			return []byte("nil should not be recomputed yet"), ttl, nil
		})
		require.NoError(t, err)
		assert.Nil(t, v)
	})
}
