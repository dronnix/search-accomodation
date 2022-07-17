package storage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMigrate(t *testing.T) {
	ctx := context.Background()
	pool, teardown := testConnectionPool(ctx, t)
	defer teardown()
	ver, err := Migrate(ctx, pool, "testdata")
	require.NoError(t, err)
	assert.Equal(t, int32(2), ver)
}

// TODO: Add tests for negative cases.
