package geolocation

import "net"

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

func NewIPLocationFromStrings(ip, countryCode, countryName, city, latitude, longitude, mystery string) (IPLocation, error) {
	return IPLocation{}, nil
}
