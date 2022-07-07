package geolocation

import "net"

// IPObservation - IP address was observed in some location.Can be obtained from CSV or other sources.
type IPObservation struct {
	IP           net.IP
	CountryCode  string
	Country      string
	City         string
	Coordinate   Coordinate
	MysteryValue uint64
}

// Coordinate represents WGS84 coordinate.
type Coordinate struct {
	Lat float64
	Lon float64
}
