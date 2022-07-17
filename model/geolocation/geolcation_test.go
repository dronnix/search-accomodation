package geolocation_test

import (
	"net"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dronnix/search-accomodation/model/geolocation"
)

func TestNewIPLocationFromStrings(t *testing.T) {
	t.Parallel()
	type args struct {
		ip          string
		countryCode string
		countryName string
		city        string
		latitude    string
		longitude   string
		mystery     string
	}
	tests := []struct {
		name    string
		args    args
		want    geolocation.IPLocation
		wantErr bool
	}{
		{
			name: "all data is valid",
			args: args{
				ip:          "8.8.8.8",
				countryCode: "UK",
				countryName: "United Kingdom",
				city:        "London",
				latitude:    "1.23",
				longitude:   "-0.42",
				mystery:     "42",
			},
			want: geolocation.IPLocation{
				IP:          net.IPv4(8, 8, 8, 8),
				CountryCode: "UK",
				CountryName: "United Kingdom",
				City:        "London",
				Coordinate: geolocation.Coordinate{
					Lat: 1.23,
					Lon: -0.42,
				},
				MysteryValue: 42,
			},
			wantErr: false,
		},
		{
			name: "non-valid IP address",
			args: args{
				ip:          "8.8.8.X",
				countryCode: "UK",
				countryName: "United Kingdom",
				city:        "London",
				latitude:    "1.23",
				longitude:   "-0.42",
				mystery:     "42",
			},
			want:    geolocation.IPLocation{},
			wantErr: true,
		},
		{
			name: "non-valid country code",
			args: args{
				ip:          "8.8.8.8",
				countryCode: "",
				countryName: "United Kingdom",
				city:        "London",
				latitude:    "1.23",
				longitude:   "-0.42",
				mystery:     "42",
			},
			want:    geolocation.IPLocation{},
			wantErr: true,
		},
		{
			name: "non-valid country name",
			args: args{
				ip:          "8.8.8.8",
				countryCode: "UK",
				countryName: "No",
				city:        "London",
				latitude:    "1.23",
				longitude:   "-0.42",
				mystery:     "42",
			},
			want:    geolocation.IPLocation{},
			wantErr: true,
		},
		{
			name: "non-valid city",
			args: args{
				ip:          "8.8.8.8",
				countryCode: "UK",
				countryName: "United Kingdom",
				city:        "",
				latitude:    "1.23",
				longitude:   "-0.42",
				mystery:     "42",
			},
			want:    geolocation.IPLocation{},
			wantErr: true,
		},
		{
			name: "non-valid latitude",
			args: args{
				ip:          "8.8.8.8",
				countryCode: "UK",
				countryName: "United Kingdom",
				city:        "London",
				latitude:    "95.1",
				longitude:   "-0.42",
				mystery:     "42",
			},
			want:    geolocation.IPLocation{},
			wantErr: true,
		},
		{
			name: "non-valid longitude",
			args: args{
				ip:          "8.8.8.8",
				countryCode: "UK",
				countryName: "United Kingdom",
				city:        "London",
				latitude:    "1.23",
				longitude:   "",
				mystery:     "42",
			},
			want:    geolocation.IPLocation{},
			wantErr: true,
		},
		{
			name: "non-valid mystery value",
			args: args{
				ip:          "8.8.8.8",
				countryCode: "UK",
				countryName: "United Kingdom",
				city:        "London",
				latitude:    "1.23",
				longitude:   "-0.42",
				mystery:     "real mystery",
			},
			want:    geolocation.IPLocation{},
			wantErr: true,
		},
		// TODO: Add test cases according to coverage map.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := geolocation.NewIPLocationFromStrings(tt.args.ip, tt.args.countryCode, tt.args.countryName,
				tt.args.city, tt.args.latitude, tt.args.longitude, tt.args.mystery)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewIPLocationFromStrings() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewIPLocationFromStrings() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIPLocation_MD5(t *testing.T) {
	t.Parallel()
	loc := geolocation.IPLocation{
		IP:          net.IPv4(8, 8, 8, 8),
		CountryCode: "UK",
		CountryName: "United Kingdom",
		City:        "London",
		Coordinate: geolocation.Coordinate{
			Lat: 1.23,
			Lon: -0.42,
		},
		MysteryValue: 2342,
	}
	notEqualLoc := geolocation.IPLocation{
		IP:          net.IPv4(8, 8, 8, 8),
		CountryCode: "UK",
		CountryName: "United Kingdom",
		City:        "London",
		Coordinate: geolocation.Coordinate{
			Lat: 1.23,
			Lon: -0.42,
		},
		MysteryValue: 2341, // <-- DIFFERENT
	}
	equalLoc := geolocation.IPLocation{
		IP:          net.IPv4(8, 8, 8, 8),
		CountryCode: "UK",
		CountryName: "United Kingdom",
		City:        "London",
		Coordinate: geolocation.Coordinate{
			Lat: 1.23,
			Lon: -0.42,
		},
		MysteryValue: 2342,
	}
	assert.Equal(t, loc.MD5(), equalLoc.MD5())
	assert.NotEqual(t, loc.MD5(), notEqualLoc.MD5())
}

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
