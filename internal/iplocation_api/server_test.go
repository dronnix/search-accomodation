package iplocation_api

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/dronnix/search-accomodation/api"
	"github.com/dronnix/search-accomodation/model/geolocation"
)

func Test_ipLocationServer_GetV1Iplocation_OK(t *testing.T) {
	t.Parallel()
	fetcher := new(fetcherMock)
	fetcher.On("FetchLocationsByIP", mock.Anything, mock.Anything).Return(locations[:1], nil).Once()
	server := NewIpLocationServer(fetcher)
	req := httptest.NewRequest(http.MethodGet, "/v1/iplocation?ip=1.2.3.4", nil)
	w := httptest.NewRecorder()
	handler := api.Handler(server)
	handler.ServeHTTP(w, req)
	res := w.Result()
	defer res.Body.Close()
	require.Equal(t, http.StatusOK, res.StatusCode)
	body, _ := io.ReadAll(res.Body)
	const expected = "{\"city\":\"London\",\"country\":\"United Kingdom\",\"country_code\":\"UK\",\"latitude\":51.5,\"longitude\":-0.1}" //nolint:lll
	require.Equal(t, expected, string(body))
}

func Test_ipLocationServer_GetV1Iplocation_Ambiguous(t *testing.T) {
	t.Parallel()
	fetcher := new(fetcherMock)
	fetcher.On("FetchLocationsByIP", mock.Anything, mock.Anything).Return(locations, nil).Once()
	server := NewIpLocationServer(fetcher)
	req := httptest.NewRequest(http.MethodGet, "/v1/iplocation?ip=1.2.3.4", nil)
	w := httptest.NewRecorder()
	handler := api.Handler(server)
	handler.ServeHTTP(w, req)
	res := w.Result()
	defer res.Body.Close()
	require.Equal(t, http.StatusNoContent, res.StatusCode)
}

type fetcherMock struct {
	mock.Mock
}

func (f *fetcherMock) FetchLocationsByIP(ctx context.Context, ip net.IP) ([]geolocation.IPLocation, error) {
	args := f.Called(ctx, ip)
	return args.Get(0).([]geolocation.IPLocation), args.Error(1) //nolint:wrapcheck
}

var locations = []geolocation.IPLocation{
	{
		IP:          net.IPv4(1, 2, 3, 4),
		CountryCode: "UK",
		CountryName: "United Kingdom",
		City:        "London",
		Coordinate: geolocation.Coordinate{
			Lat: 51.5,
			Lon: -0.1,
		},
		MysteryValue: 42,
	},
	{
		IP:          net.IPv4(1, 2, 3, 5),
		CountryCode: "UK",
		CountryName: "United Kingdom",
		City:        "London",
		Coordinate: geolocation.Coordinate{
			Lat: 51.5,
			Lon: -0.1,
		},
		MysteryValue: 42,
	},
}
