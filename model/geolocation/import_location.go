package geolocation

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"time"
)

// ImportIPLocations - imports IP locations with providing statistics.
// Returns a wrapped error if any problem occurs.
func ImportIPLocations(
	ctx context.Context,
	importer IPLocationImporter,
	storer IPLocationStorer,
) (ImportStatistics, error) {
	const batchSize = 65536
	totalStats := ImportStatistics{}
	depup := make(ipLocationsDeduplicator)
	start := time.Now()

	for {
		ipLocations, stats, err := importer.ImportNextBatch(ctx, batchSize)
		if err != nil {
			if errors.Is(err, io.EOF) {
				totalStats.TimeSpent = time.Since(start)
				return totalStats, nil
			}
			return ImportStatistics{}, fmt.Errorf("failed to import ip locations: %w", err)
		}

		var dups int
		ipLocations, dups = depup.deduplicate(ipLocations)
		stats.ApplyDuplicates(dups)

		if err = storer.StoreIPLocations(ctx, ipLocations); err != nil {
			return ImportStatistics{}, fmt.Errorf("failed to store ip locations: %w", err)
		}
		totalStats.Add(stats)
	}
}

// IPLocationImporter - interface for importing IP locations from some source.
type IPLocationImporter interface {
	// ImportNextBatch returns io.EOF when no more data is available.
	ImportNextBatch(ctx context.Context, size int) ([]IPLocation, ImportStatistics, error)
}

// IPLocationStorer - interface for storing IP locations.
type IPLocationStorer interface {
	StoreIPLocations(ctx context.Context, locations []IPLocation) error
}

// ImportStatistics - provides statistics about import process.
type ImportStatistics struct {
	Imported   int
	NonValid   int
	Duplicated int
	TimeSpent  time.Duration
}

func (s *ImportStatistics) Add(other ImportStatistics) {
	s.Imported += other.Imported
	s.NonValid += other.NonValid
	s.Duplicated += other.Duplicated
}

func (s *ImportStatistics) ApplyDuplicates(dups int) {
	s.Duplicated += dups
	s.Imported -= dups
}

func (s *ImportStatistics) Total() int {
	return s.Imported + s.NonValid + s.Duplicated
}

type ipLocationsDeduplicator map[[md5.Size]byte]bool

func (d *ipLocationsDeduplicator) deduplicate(locations []IPLocation) (result []IPLocation, duplicated int) {
	for i := 0; i < len(locations); i++ {
		if _, ok := (*d)[locations[i].MD5()]; ok {
			locations[i] = locations[len(locations)-1]
			locations = locations[:len(locations)-1]
			duplicated++
			continue
		}
		(*d)[locations[i].MD5()] = true
	}
	return locations, duplicated
}
