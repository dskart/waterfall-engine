package store

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func NewTestStore(ctx context.Context, t *testing.T) *Store {
	s, err := New(ctx, Config{
		InMemory: true,
	})
	require.NoError(t, err)
	return s
}
