package main

import (
	"context"
	"fmt"
	"os"

	"github.com/dronnix/search-accomodation/domain/geolocation"
	"github.com/dronnix/search-accomodation/internal/pkg/iplocation_importer"
	"github.com/dronnix/search-accomodation/storage"
)

func main() {
	f, err := os.Open("data_dump.csv")
	if err != nil {
		panic(err)
	}
	importer, err := iplocation_importer.NewCSVImporter(f)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	pool, err := storage.CreateConnectionPool(ctx, "postgres://test:test@localhost:5432/test")
	if err != nil {
		panic(err)
	}
	storer := storage.NewIPLocationStorage(pool)
	if err := storer.MigrateUp(ctx, "storage/migrations/iplocation"); err != nil {
		panic(err)
	}
	stats, err := geolocation.ImportIPLocations(ctx, importer, storer)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", stats)
}
