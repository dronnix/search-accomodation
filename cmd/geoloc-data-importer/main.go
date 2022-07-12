package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/dronnix/search-accomodation/internal/pkg/iplocation_importer"
)

func main() {

	f, err := os.Open("data_dump.csv")
	if err != nil {
		panic(err)
	}
	importer, err := iplocation_importer.NewCSVImporter(f)
	if err != nil {
		panic(err)
	}

	places := make(map[string]map[string]struct{})
	for {
		locations, _, err := importer.ImportNextBatch(context.Background(), 4096)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			panic(err)
		}
		for _, location := range locations {
			if places[location.CountryCode] == nil {
				places[location.CountryCode] = make(map[string]struct{})
			}
			places[location.CountryCode][location.CountryName] = struct{}{}
		}
	}
	fmt.Println(len(places))
	for countryCode, country := range places {
		fmt.Println(countryCode)
		for countryName := range country {
			fmt.Println("\t", countryName)
		}
	}
}
