package storage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func setupIPLocationStorage(ctx context.Context, t *testing.T) (storage *IPLocationStorage, teardown func()) {
	pool, teardown := testConnectionPool(ctx, t)
	return NewIPLocationStorage(pool), teardown
}

func setUpDB(t *testing.T) (context.Context, *IPLocationStorage, func()) {
	ctx := context.Background()
	storage, teardown := setupIPLocationStorage(ctx, t)
	return ctx, storage, teardown
}

// Test utilities
func TestStorage_StoreObservations_MigrateUp(t *testing.T) {
	ctx := context.Background()
	storage, teardown := setupIPLocationStorage(ctx, t)
	defer teardown()
	require.NoError(t, storage.MigrateUp(ctx, migrationsDir))
}

const migrationsDir = "migrations/iplocation"
