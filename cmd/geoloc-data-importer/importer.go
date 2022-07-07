package main

import (
	"time"

	"github.com/dronnix/search-accomodation/model/geolocation"
)

type csvImporter struct{}

func newCsvImporter(pathToCSVFile string) *csvImporter {

	return &csvImporter{}
}

type csvImportStats struct {
	TotalRecords    int
	StoredRecords   int
	NonValidRecords int
	TimeSpent       time.Duration
}

func (c *csvImporter) ImportNextBatch(size int) ([]geolocation.IPObservation, error) {
	return nil, nil
}
