package geolocation

import (
	"fmt"
	"strconv"
)

// Coordinate represents WGS84 coordinate.
type Coordinate struct { // TODO: Replace with geo-library.
	Lat float64
	Lon float64
}

// NewCoordinateFromStrings - creates Coordinate from strings representation.
func NewCoordinateFromStrings(sLat, sLon string) (Coordinate, error) {
	lat, err := strconv.ParseFloat(sLat, 64)
	if err != nil {
		return Coordinate{}, fmt.Errorf("failed to parse latitude: %w", err)
	}
	lon, err := strconv.ParseFloat(sLon, 64)
	if err != nil {
		return Coordinate{}, fmt.Errorf("failed to parse longitude: %w", err)
	}
	coord := Coordinate{Lat: lat, Lon: lon}
	if err = coord.Validate(); err != nil {
		return Coordinate{}, fmt.Errorf("coordinate is invalid: %w", err)
	}
	return coord, nil
}

func (c Coordinate) Validate() error {
	if c.Lat < -(90+epsilon) || c.Lat > (90+epsilon) {
		return fmt.Errorf("latitude is out of bounds [-90;90]: %f", c.Lat)
	}
	if c.Lon < -(180+epsilon) || c.Lon > (180+epsilon) {
		return fmt.Errorf("longitude is out of bounds [-180;180]: %f", c.Lon)
	}
	return nil
}

const epsilon = 0.000001
