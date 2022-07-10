package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"

	"github.com/dronnix/search-accomodation/domain/geolocation"
)

type csvImporter struct {
	csvReader *csv.Reader
}

func newCSVImporter(r io.Reader) (*csvImporter, error) {
	csvReader := csv.NewReader(r)

	header, err := csvReader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read csv header: %w", err)
	}
	if len(header) != 7 {
		return nil, fmt.Errorf("header must contain 7 columns: %v", header)
	}
	// TODO: validate header fields.
	return &csvImporter{csvReader: csvReader}, nil
}

type ImportStats struct {
	Imported int
	NonValid int
}

func (c *csvImporter) ImportNextBatch(size int) ([]geolocation.IPLocation, ImportStats, error) {
	ipLocations := make([]geolocation.IPLocation, 0, size)
	stats := ImportStats{}
	for i := 0; i < size; i++ {
		rec, err := c.csvReader.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				if i == 0 {
					return nil, ImportStats{}, io.EOF
				}
				break
			}
			return nil, ImportStats{}, fmt.Errorf("failed to read record: %w", err)
		}
		if len(rec) != 7 {
			stats.NonValid++
			continue
		}

		location, err := geolocation.NewIPLocationFromStrings(rec[0], rec[1], rec[2], rec[3], rec[4], rec[5], rec[6])
		if err != nil {
			stats.NonValid++
			continue
		}
		ipLocations = append(ipLocations, location)
		stats.Imported++
	}
	return ipLocations, stats, nil
}

func (s *ImportStats) Add(other ImportStats) {
	s.Imported += other.Imported
	s.NonValid += other.NonValid
}
