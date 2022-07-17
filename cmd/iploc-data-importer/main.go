package main

import (
	"context"
	"fmt"
	"os"

	"github.com/dronnix/search-accomodation/internal/flags"
	"github.com/dronnix/search-accomodation/internal/iplocation_importer"
	"github.com/dronnix/search-accomodation/model/geolocation"
	"github.com/dronnix/search-accomodation/storage"
)

type options struct {
	PathToCSV string `long:"path-to-csv" default:"data_dump.csv" env:"PATH_TO_CSV"`
	*flags.Postgres
}

const exitCodeOK = 0
const exitCodeError = 1

func main() {
	os.Exit(_main())
}

func _main() int {
	opts := &options{}
	flags.Parse(opts)

	importer, err := setupImporter(opts.PathToCSV)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not setup importer: %v\n", err)
		return exitCodeError
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	storage, err := setupStorage(ctx, opts.Postgres)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not setup storage: %v\n", err)
		return exitCodeError
	}

	stats, err := geolocation.ImportIPLocations(ctx, importer, storage)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not import IP locations: %v\n", err)
		return exitCodeError
	}

	fmt.Printf("Time spent(sec): %d\n", int(stats.TimeSpent.Seconds()))
	fmt.Printf("Total records found: %d\n", stats.Total())
	fmt.Printf("Non-valid records: %d\n", stats.NonValid)
	fmt.Printf("Duplicated records: %d\n", stats.Duplicated)
	fmt.Printf("Imported records: %d\n", stats.Imported)

	return exitCodeOK
}

func setupImporter(pathToCSV string) (*iplocation_importer.CSVImporter, error) {
	f, err := os.Open(pathToCSV)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %w", err)
	}
	return iplocation_importer.NewCSVImporter(f) //nolint:wrapcheck
}

func setupStorage(ctx context.Context, opts *flags.Postgres) (*storage.IPLocationStorage, error) {
	pool, err := storage.CreateConnectionPool(ctx, opts.PostgresConnectionString())
	if err != nil {
		return nil, fmt.Errorf("could not create connection pool: %w", err)
	}
	s := storage.NewIPLocationStorage(pool)
	const migrationsPath = "storage/migrations/iplocation" // TODO: move to config
	if err := s.MigrateUp(ctx, migrationsPath); err != nil {
		return nil, fmt.Errorf("could not migrate up: %w", err)
	}
	return s, nil
}
