package geolocation_test

import (
	"context"
	"io"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/dronnix/search-accomodation/model/geolocation"
)

func TestImportStatistics_Add(t *testing.T) {
	t.Parallel()
	s1 := geolocation.ImportStatistics{
		Imported:   7,
		NonValid:   3,
		Duplicated: 1,
	}
	s2 := geolocation.ImportStatistics{
		Imported:   3,
		NonValid:   2,
		Duplicated: 1,
	}
	s1.Add(s2)
	require.Equal(t, geolocation.ImportStatistics{
		Imported:   10,
		NonValid:   5,
		Duplicated: 2,
	}, s1)
}

func TestImportStatistics_ApplyDuplicates(t *testing.T) {
	t.Parallel()
	s1 := geolocation.ImportStatistics{
		Imported:   7,
		NonValid:   3,
		Duplicated: 1,
	}
	s1.ApplyDuplicates(3)
	require.Equal(t, geolocation.ImportStatistics{
		Imported:   4,
		NonValid:   3,
		Duplicated: 4,
	}, s1)
}

func TestImportStatistics_Total(t *testing.T) {
	t.Parallel()
	s1 := geolocation.ImportStatistics{
		Imported:   7,
		NonValid:   3,
		Duplicated: 1,
	}
	require.Equal(t, 11, s1.Total())
}

func TestImportIPLocations(t *testing.T) {
	t.Parallel()
	importer, storer := new(importerMock), new(storerMock)
	importer.On("ImportNextBatch", mock.Anything, mock.AnythingOfType("int")).Return(
		locations, geolocation.ImportStatistics{Imported: 3, TimeSpent: 0}, nil).Once()
	importer.On("ImportNextBatch", mock.Anything, mock.AnythingOfType("int")).Return(
		[]geolocation.IPLocation{}, geolocation.ImportStatistics{}, io.EOF).Once()

	storer.On("StoreIPLocations", mock.Anything, locations).Return(nil).Once()

	stats, err := geolocation.ImportIPLocations(context.Background(), importer, storer)
	require.NoError(t, err)
	assert.Equal(t, 3, stats.Imported)
	assert.Equal(t, 0, stats.Duplicated)
	assert.Equal(t, 0, stats.NonValid)

	importer.AssertExpectations(t)
	storer.AssertExpectations(t)
}

// TODO: Add more cases for ImportIPLocations.

type importerMock struct {
	mock.Mock
}

func (i *importerMock) ImportNextBatch(
	ctx context.Context,
	size int,
) ([]geolocation.IPLocation, geolocation.ImportStatistics, error) {
	args := i.Called(ctx, size)
	//nolint:wrapcheck
	return args.Get(0).([]geolocation.IPLocation), args.Get(1).(geolocation.ImportStatistics), args.Error(2)
}

type storerMock struct {
	mock.Mock
}

func (s *storerMock) StoreIPLocations(ctx context.Context, locations []geolocation.IPLocation) error {
	args := s.Called(ctx, locations)
	return args.Error(0) //nolint:wrapcheck
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
		CountryCode: "FR",
		CountryName: "France",
		City:        "Paris",
		Coordinate: geolocation.Coordinate{
			Lat: 51.5,
			Lon: -0.1,
		},
		MysteryValue: 23,
	},
	{
		IP:          net.IPv4(1, 2, 3, 6),
		CountryCode: "NZ",
		CountryName: "New Zealand",
		City:        "Auckland",
		Coordinate: geolocation.Coordinate{
			Lat: -36.8,
			Lon: 174.7,
		},
		MysteryValue: 42,
	},
}
