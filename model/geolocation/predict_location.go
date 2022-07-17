package geolocation

import (
	"context"
	"errors"
	"fmt"
	"net"
)

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

// IPLocationFetcher - interface for fetching locations by IP address.
type IPLocationFetcher interface {
	// FetchLocationsByIP returns all possible locations for given IP address.
	// If no locations are found, returns empty slice.
	FetchLocationsByIP(ctx context.Context, ip net.IP) ([]IPLocation, error)
}
