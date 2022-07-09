package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/dronnix/search-accomodation/domain/geolocation"
)

type csvImporter struct {
	reader io.Reader
}

func newCsvImporter(r io.Reader) *csvImporter {
	return &csvImporter{reader: r}
}

type csvImportStats struct {
	TotalRecords    int
	StoredRecords   int
	NonValidRecords int
	TimeSpent       time.Duration
}

func (c *csvImporter) ImportNextBatch(size int) ([]geolocation.IPLocation, error) {
	csvReader := csv.NewReader(c.reader)

	total, properLen := 0, 0
	ipToRecord := make(map[string]int)

	for i := 0; i < size; i++ {
		rec, err := csvReader.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, fmt.Errorf("failed to read record: %w", err)
		}
		total++
		if len(rec) < 7 {
			continue
		}
		properLen++
		ipToRecord[rec[0]]++
	}
	fmt.Printf("Total records: %d, proper records: %d\n", total, properLen)
	fmt.Printf("IPs: %d\n", len(ipToRecord))
	for k, v := range ipToRecord {
		if v > 1 {
			fmt.Printf("%s: %d\n", k, v)
		}
	}
	return nil, nil
}
