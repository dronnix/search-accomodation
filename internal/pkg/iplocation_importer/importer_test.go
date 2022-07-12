package iplocation_importer_test

import (
	"context"
	"errors"
	"io"
	"net"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dronnix/search-accomodation/domain/geolocation"
	"github.com/dronnix/search-accomodation/internal/pkg/iplocation_importer"
)

func getExampleFileReader(t *testing.T) (io.Reader, func()) {
	const path = "../../../data_dump.csv"
	f, err := os.Open(path)
	require.NoError(t, err)
	return f, func() { f.Name() }
}

func Test_csvImporter_ImportNextBatch_CheckOnExampleFile(t *testing.T) {
	// TODO: Disable it, or use short version!
	t.Parallel()
	reader, cleanup := getExampleFileReader(t)
	defer cleanup()
	importer, err := iplocation_importer.NewCSVImporter(reader)
	require.NoError(t, err)
	totalStats := geolocation.ImportStatistics{}
	for {
		records, stats, err := importer.ImportNextBatch(context.Background(), 7)
		if errors.Is(err, io.EOF) {
			break
		}
		require.NoError(t, err)
		require.Equal(t, stats.Imported, len(records))
		totalStats.Add(stats)
	}
	assert.Equal(t, 899431, totalStats.Imported)
	assert.Equal(t, 100569, totalStats.NonValid)
}

func Test_csvImporter_ImportNextBatch(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name            string
		csvData         string
		sizeArg         int
		wantLocation    []geolocation.IPLocation
		wantStats       geolocation.ImportStatistics
		wantCreationErr bool
		wantImportErr   bool
	}{
		{
			name:            "wrong header",
			csvData:         "ip,country,region,city,lat,lon",
			sizeArg:         1,
			wantLocation:    nil,
			wantStats:       geolocation.ImportStatistics{},
			wantCreationErr: true,
			wantImportErr:   false,
		},
		{
			name:    "valid data",
			csvData: validHeader + validRecord,
			sizeArg: 2,
			wantLocation: []geolocation.IPLocation{
				{
					IP:           net.IPv4(200, 106, 141, 15),
					CountryCode:  "SI",
					CountryName:  "Nepal",
					City:         "DuBuquemouth",
					Coordinate:   geolocation.Coordinate{Lat: -84.87503094689836, Lon: 7.206435933364332},
					MysteryValue: 7823011346,
				},
			},
			wantStats:       geolocation.ImportStatistics{Imported: 1, NonValid: 0},
			wantCreationErr: false,
			wantImportErr:   false,
		},
		{
			name:    "invalid ip",
			csvData: validHeader + validRecord + invalidRecord,
			sizeArg: 7,
			wantLocation: []geolocation.IPLocation{
				{ // SI,Nepal,DuBuquemouth,-84.87503094689836,7.206435933364332,7823011346
					IP:           net.IPv4(200, 106, 141, 15),
					CountryCode:  "SI",
					CountryName:  "Nepal",
					City:         "DuBuquemouth",
					Coordinate:   geolocation.Coordinate{Lat: -84.87503094689836, Lon: 7.206435933364332},
					MysteryValue: 7823011346,
				},
			},
			wantStats:       geolocation.ImportStatistics{Imported: 1, NonValid: 1},
			wantCreationErr: false,
			wantImportErr:   false,
		},
		// TODO: Add more tests according to the coverage map.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			t.Parallel()
			importer, err := iplocation_importer.NewCSVImporter(strings.NewReader(tt.csvData))
			if tt.wantCreationErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			loc, stats, err := importer.ImportNextBatch(context.Background(), tt.sizeArg)
			if tt.wantImportErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			assert.Equal(t, tt.wantLocation, loc)
			assert.Equal(t, tt.wantStats, stats)
		})
	}
}

const validHeader = "ip_address,country_code,country,city,latitude,longitude,mystery_value\n"
const validRecord = "200.106.141.15,SI,Nepal,DuBuquemouth,-84.87503094689836,7.206435933364332,7823011346\n"
const invalidRecord = "XXX.103.7.140,CZ,Nicaragua,New Neva,-68.31023296602508,-37.62435199624531,7301823115\n"
