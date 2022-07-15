package geolocation

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"net"
	"regexp"
	"strconv"
)

// IPLocation - IP address was observed in some location.Can be obtained from CSV or other sources.
type IPLocation struct {
	IP          net.IP
	CountryCode string
	CountryName string
	City        string
	Coordinate
	MysteryValue uint64
}

// ImportIPLocations - imports IP locations with providing statistics.
// Returns a wrapped error if any problem occurs.
func ImportIPLocations(
	ctx context.Context,
	importer IPLocationImporter,
	storer IPLocationStorer) (ImportStatistics, error) {
	const batchSize = 4096
	totalStats := ImportStatistics{}
	depup := make(ipLocationsDeduplicator)
	// TODO: Measure time spent on import.
	for {
		ipLocations, stats, err := importer.ImportNextBatch(ctx, batchSize)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return totalStats, nil
			}
			return stats, fmt.Errorf("failed to import ip locations: %w", err)
		}

		var dups int
		ipLocations, dups = depup.deduplicate(ipLocations)
		stats.ApplyDuplicates(dups)

		if err = storer.StoreIPLocations(ctx, ipLocations); err != nil {
			return stats, fmt.Errorf("failed to store ip locations: %w", err)
		}
		totalStats.Add(stats)
	}
}

var ErrIPLocationNotFound = errors.New("ip location not found")
var ErrIPLocationAmbiguous = errors.New("ip location is ambiguous")

// PredictIPLocation figures out IP location from IP address, using given fetcher.
// Returns ErrIPLocationNotFound if IP location is not found.
// Returns ErrIPLocationAmbiguous if IP more than one location known for the IP.
// Returns a wrapped error if any other error occurs.
func PredictIPLocation(ctx context.Context, ip net.IP, fetcher IPLocationFetcher) (IPLocation, error) {
	locations, err := fetcher.FetchLocationsByIP(ctx, ip)
	if err != nil {
		return IPLocation{}, fmt.Errorf("failed to fetch ip locations: %w", err)
	}
	if len(locations) == 0 {
		return IPLocation{}, ErrIPLocationNotFound
	}
	if len(locations) > 1 {
		return IPLocation{}, ErrIPLocationAmbiguous
	}
	return locations[0], nil
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

// IPLocationFetcher - interface for fetching locations by IP address.
type IPLocationFetcher interface {
	// FetchLocationsByIP returns all possible locations for given IP address.
	// If no locations are found, returns empty slice.
	FetchLocationsByIP(ctx context.Context, ip net.IP) ([]IPLocation, error)
}

// ImportStatistics - provides statistics about import process.
type ImportStatistics struct {
	Imported   int
	NonValid   int
	Duplicated int
}

// NewIPLocationFromStrings - creates IPLocation from strings representation. Useful for CSVs, logs, etc.
func NewIPLocationFromStrings(
	ip,
	countryCode,
	countryName,
	city,
	latitude,
	longitude,
	mystery string,
) (IPLocation, error) {
	ipAddr := net.ParseIP(ip)
	if ipAddr == nil {
		return IPLocation{}, fmt.Errorf("failed to parse ip address: %s", ip)
	}

	// TODO: Is better to validate names through countries/cities catalog with normalizing.
	if !validCountryCode(countryCode) {
		return IPLocation{}, fmt.Errorf("country code must be 2 characters: %s", countryCode)
	}
	if len(countryName) < 4 {
		return IPLocation{}, fmt.Errorf("country name must be at least 4 characters: %s", countryName)
	}
	if city == "" {
		return IPLocation{}, fmt.Errorf("city name is empty")
	}

	coord, err := NewCoordinateFromStrings(latitude, longitude)
	if err != nil {
		return IPLocation{}, fmt.Errorf("failed to parse coordinate: %w", err)
	}

	mysteryValue, err := strconv.ParseUint(mystery, 10, 64)
	if err != nil {
		return IPLocation{}, fmt.Errorf("failed to parse mystery value: %w", err)
	}

	return IPLocation{
		IP:           ipAddr,
		CountryCode:  countryCode,
		CountryName:  countryName,
		City:         city,
		Coordinate:   coord,
		MysteryValue: mysteryValue,
	}, nil
}

func (l *IPLocation) MD5() [md5.Size]byte {
	var b bytes.Buffer
	_ = gob.NewEncoder(&b).Encode(*l)
	return md5.Sum(b.Bytes())
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

var validCountryCode = regexp.MustCompile(`^[a-zA-Z]{2}$`).MatchString

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
