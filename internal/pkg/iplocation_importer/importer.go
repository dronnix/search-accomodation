package iplocation_importer

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"

	"github.com/dronnix/search-accomodation/domain/geolocation"
)

// CSVImporter imports ip locations from CSV files with following structure:
// ip_address,country_code,country,city,latitude,longitude,mystery_value
type CSVImporter struct {
	csvReader *csv.Reader
}

// NewCSVImporter creates a CSVImporter from reader
func NewCSVImporter(r io.Reader) (*CSVImporter, error) {
	csvReader := csv.NewReader(r)

	header, err := csvReader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read csv header: %w", err)
	}
	if len(header) != 7 {
		return nil, fmt.Errorf("header must contain 7 columns: %v", header)
	}
	// TODO: validate header fields.
	return &CSVImporter{csvReader: csvReader}, nil
}

func (c *CSVImporter) ImportNextBatch(
	ctx context.Context,
	size int,
) ([]geolocation.IPLocation, geolocation.ImportStatistics, error) {
	ipLocations := make([]geolocation.IPLocation, 0, size)
	stats := geolocation.ImportStatistics{}
	for i := 0; i < size; i++ {
		select {
		case <-ctx.Done():
			return ipLocations, stats, ctx.Err() //nolint: wrapcheck
		default:
		}
		rec, err := c.csvReader.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				if i == 0 {
					return nil, geolocation.ImportStatistics{}, io.EOF
				}
				break
			}
			return nil, geolocation.ImportStatistics{}, fmt.Errorf("failed to read record: %w", err)
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
