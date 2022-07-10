package geolocation

import (
	"fmt"
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

// Coordinate represents WGS84 coordinate.
type Coordinate struct {
	Lat float64
	Lon float64
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

var validCountryCode = regexp.MustCompile(`^[a-zA-Z]{2}$`).MatchString

const epsilon = 0.000001
