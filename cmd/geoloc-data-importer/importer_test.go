package main

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func getExampleFileReader(t *testing.T) io.Reader {
	const path = "../../data_dump.csv"
	f, err := os.Open(path)
	require.NoError(t, err)
	return f
}

func Test_csvImporter_ImportNextBatch_CheckOnExampleFile(t *testing.T) {
	importer := newCsvImporter(getExampleFileReader(t))
	_, err := importer.ImportNextBatch(1024 * 1024)
	require.NoError(t, err)
}
