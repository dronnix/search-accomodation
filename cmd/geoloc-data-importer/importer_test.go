package main

import (
	"errors"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getExampleFileReader(t *testing.T) (io.Reader, func()) {
	const path = "../../data_dump.csv"
	f, err := os.Open(path)
	require.NoError(t, err)
	return f, func() { f.Name() }
}

func Test_csvImporter_ImportNextBatch_CheckOnExampleFile(t *testing.T) {
	t.Parallel()
	reader, cleanup := getExampleFileReader(t)
	defer cleanup()
	importer, err := newCSVImporter(reader)
	require.NoError(t, err)
	totalStats := ImportStats{}
	for {
		records, stats, err := importer.ImportNextBatch(7)
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
