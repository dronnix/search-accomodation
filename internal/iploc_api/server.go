package iploc_api

import (
	"encoding/json"
	"errors"
	"net"
	"net/http"

	"github.com/dronnix/search-accomodation/api"
	"github.com/dronnix/search-accomodation/model/geolocation"
)

// IPLocationServer is handler-implementation for auto-generated API stub.
type IPLocationServer struct {
	fetcher geolocation.IPLocationFetcher
}

func NewIpLocationServer(fetcher geolocation.IPLocationFetcher) *IPLocationServer {
	return &IPLocationServer{fetcher: fetcher}
}

// GetV1Iplocation is handler-implementation for auto-generated API stub.
func (s *IPLocationServer) GetV1Iplocation(w http.ResponseWriter, r *http.Request, params api.GetV1IplocationParams) {
	ip := net.ParseIP(params.Ip)
	if ip == nil {
		s.sendResponse(http.StatusBadRequest, w, api.Error{ErrorDetails: "Invalid IP address"})
		return
	}

	location, err := geolocation.PredictIPLocation(r.Context(), ip, s.fetcher)
	if err != nil {
		if errors.Is(err, geolocation.ErrIPLocationNotFound) || errors.Is(err, geolocation.ErrIPLocationAmbiguous) {
			s.sendResponse(http.StatusNoContent, w, nil)
		} else {
			// TODO: Log error, don't expose it to the user.
			s.sendResponse(http.StatusServiceUnavailable, w, api.Error{ErrorDetails: err.Error()})
		}
		return
	}

	s.sendResponse(http.StatusOK, w, api.IpLocation{
		City:        location.City,
		Country:     location.CountryName,
		CountryCode: location.CountryCode,
		Latitude:    location.Lat,
		Longitude:   location.Lon,
	})
}

func (s *IPLocationServer) sendResponse(code int, w http.ResponseWriter, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		body, err := json.Marshal(data)
		if err != nil {
			panic(err) // Exceptional situation - response structure must be marshalable.
		}
		_, _ = w.Write(body) // TODO:Log error.
	}
}
