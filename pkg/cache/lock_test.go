package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLock(t *testing.T) {
	var lock Lock

	invoked := make(chan struct{})
	block := make(chan struct{})

	var invocations int

	f := func() (interface{}, error) {
		close(invoked)
		<-block
		invocations++
		return 1, nil
	}

	// This Acquire will close(invoked) then block.
	go lock.Acquire("foo", f)
	<-invoked

	// Additional calls to Acquire should return immediately.
	results := make([]*LockResult, 10)
	for i := range results {
		results[i] = lock.Acquire("foo", f)
	}

	// Unblock the lock owner.
	close(block)

	for _, result := range results {
		value, err := result.Wait()
		assert.Equal(t, 1, value)
		assert.NoError(t, err)
	}

	assert.Equal(t, invocations, 1)
}
