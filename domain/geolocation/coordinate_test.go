package geolocation_test

import (
	"testing"

	"github.com/dronnix/search-accomodation/domain/geolocation"
)

func TestCoordinates_Validate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		point   geolocation.Coordinate
		wantErr bool
	}{
		{
			name:    "should pass with Greenwich",
			point:   geolocation.Coordinate{},
			wantErr: false,
		},
		{
			name: "should pass with max values",
			point: geolocation.Coordinate{
				Lat: 90,
				Lon: 180,
			},
			wantErr: false,
		},
		{
			name: "should pass with min values",
			point: geolocation.Coordinate{
				Lat: -90,
				Lon: -180,
			},
			wantErr: false,
		},
		{
			name: "should fail with Lat overflow",
			point: geolocation.Coordinate{
				Lat: 90.01,
				Lon: 180,
			},
			wantErr: true,
		},
		{
			name: "should fail with Lat underflow",
			point: geolocation.Coordinate{
				Lat: -90.01,
				Lon: 180,
			},
			wantErr: true,
		},
		{
			name: "should fail with Lon overflow",
			point: geolocation.Coordinate{
				Lat: 90,
				Lon: 180.01,
			},
			wantErr: true,
		},
		{
			name: "should fail with Lat underflow",
			point: geolocation.Coordinate{
				Lat: 90,
				Lon: -180.01,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests { //nolint:paralleltest
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.point.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
