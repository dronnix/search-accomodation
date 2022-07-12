package storage

import (
	"context"
	"fmt"
	"net"

	"github.com/jackc/pgx/v4"
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
	table := []string{"geolocation", "ip_location"}
	columns := []string{"ip_address", "country_code", "country_name", "city", "latitude", "longitude",
		"mystery_value"}

	locs := make([][]interface{}, len(locations))
	for i := range locations {
		locs[i] = []interface{}{
			locations[i].IP,
			locations[i].CountryCode,
			locations[i].CountryName,
			locations[i].City,
			locations[i].Lat,
			locations[i].Lon,
			locations[i].MysteryValue,
		}
	}

	n, err := s.pool.CopyFrom(ctx, table, columns, pgx.CopyFromRows(locs))
	if err != nil {
		return fmt.Errorf("unable to copy observations to db: %w", err)
	}
	if n != int64(len(locations)) {
		return fmt.Errorf("stored unexpected number of observations")
	}
	return nil
}

func (s *IPLocationStorage) FetchLocationsByIP(ip net.IP) ([]geolocation.IPLocation, error) {
	//TODO implement me
	panic("implement me")
}
