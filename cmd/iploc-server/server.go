package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/dronnix/search-accomodation/api"
)

type server struct {
}

func (s *server) GetV1Iplocation(w http.ResponseWriter, r *http.Request, params api.GetV1IplocationParams) {
	ip := net.ParseIP(params.Ip)
	if ip == nil {
		w.WriteHeader(http.StatusBadRequest)
		body, err := json.Marshal(&api.Error{
			ErrorDetails: fmt.Sprintf("Invalid IP address: %s", params.Ip),
		})
		if err != nil {
			panic(err) // will be handled by recovery middleware.
		}
		w.Write(body)
		return
	}

	body, err := json.Marshal(&api.IpLocation{
		City:        "Tbilisi",
		Country:     "Georgia",
		CountryCode: "GE",
		Latitude:    23.3,
		Longitude:   42.2,
	})
	if err != nil {
		panic(err) // will be handled by recovery middleware.
	}
	w.Write(body)
}
