package geolocation

import (
	"bytes"
	"crypto/md5"
	"encoding/gob"
	"fmt"
	"net"
	"regexp"
	"strconv"
)

// IPLocation - IP address that was observed in some location.Can be obtained from CSV or other sources.
type IPLocation struct {
	IP          net.IP
	CountryCode string
	CountryName string
	City        string
	Coordinate
	MysteryValue uint64
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

var validCountryCode = regexp.MustCompile(`^[a-zA-Z]{2}$`).MatchString
