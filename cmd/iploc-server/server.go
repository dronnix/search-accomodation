package main

import (
	"encoding/json"
	"errors"
	"net"
	"net/http"

	"github.com/dronnix/search-accomodation/api"
	"github.com/dronnix/search-accomodation/model/geolocation"
)

type ipLocationServer struct {
	fetcher geolocation.IPLocationFetcher
}

func newIpLocationServer(fetcher geolocation.IPLocationFetcher) *ipLocationServer {
	return &ipLocationServer{fetcher: fetcher}
}

func (s *ipLocationServer) GetV1Iplocation(w http.ResponseWriter, r *http.Request, params api.GetV1IplocationParams) {
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

func (s *ipLocationServer) sendResponse(code int, w http.ResponseWriter, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		body, err := json.Marshal(data)
		if err != nil {
			panic(err) // Exceptional situation - response structure must be marshalable.
		}
		_, _ = w.Write(body) // TODO:Log error.
	}
}
