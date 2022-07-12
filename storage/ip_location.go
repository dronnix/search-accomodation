package storage

import (
	"context"
	"fmt"
	"net"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/dronnix/search-accomodation/domain/geolocation"
)

type IPLocationStorage struct {
	pool *pgxpool.Pool
}

var _ geolocation.IPLocationFetcher = (*IPLocationStorage)(nil)
var _ geolocation.IPLocationStorer = (*IPLocationStorage)(nil)

func NewIPLocationStorage(pool *pgxpool.Pool) *IPLocationStorage {
	return &IPLocationStorage{pool: pool}
}

func (s *IPLocationStorage) MigrateUp(ctx context.Context, migrationsDir string) error {
	version, err := Migrate(ctx, s.pool, migrationsDir)
	if err != nil {
		return fmt.Errorf("unable to migrate observations to version %d: %w", version, err)
	}
	return nil
}

func (s *IPLocationStorage) StoreIPLocations(ctx context.Context, locations []geolocation.IPLocation) error {
	//TODO implement me
	panic("implement me")
}

func (s *IPLocationStorage) FetchLocationsByIP(ip net.IP) ([]geolocation.IPLocation, error) {
	//TODO implement me
	panic("implement me")
}
