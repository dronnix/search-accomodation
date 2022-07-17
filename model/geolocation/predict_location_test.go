package geolocation_test

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/dronnix/search-accomodation/model/geolocation"
)

func TestPredictIPLocation(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	fetcher := new(fetcherMock)
	ip := net.IPv4(1, 2, 3, 6)
	fetcher.On("FetchLocationsByIP", ctx, ip).Return(locations[2:], nil).Once()

	loc, err := geolocation.PredictIPLocation(ctx, ip, fetcher)
	require.NoError(t, err)
	require.Equal(t, locations[2], loc)

	fetcher.AssertExpectations(t)
}

func TestPredictIPLocation_Multiple(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	fetcher := new(fetcherMock)
	ip := net.IPv4(1, 2, 3, 6)
	fetcher.On("FetchLocationsByIP", ctx, ip).Return(locations, nil).Once()

	_, err := geolocation.PredictIPLocation(ctx, ip, fetcher)
	require.EqualError(t, geolocation.ErrIPLocationAmbiguous, err.Error())

	fetcher.AssertExpectations(t)
}

func TestPredictIPLocation_NotFound(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	fetcher := new(fetcherMock)
	ip := net.IPv4(1, 2, 3, 6)
	fetcher.On("FetchLocationsByIP", ctx, ip).Return([]geolocation.IPLocation{}, nil).Once()

	_, err := geolocation.PredictIPLocation(ctx, ip, fetcher)
	require.EqualError(t, geolocation.ErrIPLocationNotFound, err.Error())

	fetcher.AssertExpectations(t)
}

type fetcherMock struct {
	mock.Mock
}

func (f *fetcherMock) FetchLocationsByIP(ctx context.Context, ip net.IP) ([]geolocation.IPLocation, error) {
	args := f.Called(ctx, ip)
	return args.Get(0).([]geolocation.IPLocation), args.Error(1) //nolint:wrapcheck
}
