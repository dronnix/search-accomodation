package storage

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dronnix/search-accomodation/domain/geolocation"
)

const migrationsDir = "migrations/iplocation"

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

func TestIPLocationStorage_StoreIPLocations(t *testing.T) {
	ctx, storage, teardown := setUpDB(t)
	defer teardown()
	require.NoError(t, storage.MigrateUp(ctx, migrationsDir))
	require.NoError(t, storage.StoreIPLocations(ctx, []geolocation.IPLocation{
		{
			IP:          net.IPv4(8, 8, 8, 8),
			CountryCode: "UK",
			CountryName: "United Kingdom",
			City:        "London",
			Coordinate: geolocation.Coordinate{
				Lat: 0.42,
				Lon: -0.23,
			},
			MysteryValue: 31337,
		},
	}))
	rows, err := storage.pool.Query(ctx, "SELECT COUNT(*) FROM geolocation.ip_location")
	require.NoError(t, err)
	defer rows.Close()
	count := 0
	require.True(t, rows.Next())
	require.NoError(t, rows.Scan(&count))
	assert.Equal(t, 1, count)
}
