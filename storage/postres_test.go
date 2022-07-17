package storage

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateConnectionPool_EmptyConnString(t *testing.T) {
	pool, err := CreateConnectionPool(context.Background(), "")
	require.Error(t, err)
	assert.Nil(t, pool)
}

func TestCreateConnectionPool_Default(t *testing.T) {
	pool, err := CreateConnectionPool(context.Background(), testConnectionString(testDBName))
	require.NoError(t, err)
	assert.NotNil(t, pool)
	defer pool.Close()
}

const testDBName = "test"

func testConnectionString(dbName string) string {
	host := os.Getenv("POSTGRES_HOST")
	if host == "" {
		host = "localhost"
	}
	return fmt.Sprintf("postgres://test:test@%s:5432/%s?&pool_max_conns=2", host, dbName)
}

func getTestDBName(testName string) string {
	testName = strings.ToLower(strings.ReplaceAll(testName, "/", "_"))
	const maxNameLength = 48
	testNameLength := len(testName)
	if testNameLength > maxNameLength {
		testName = testName[testNameLength-maxNameLength : testNameLength]
	}
	return fmt.Sprintf("%s_%d", testName, time.Now().UnixNano()/1000)
}

func testConnectionPool(ctx context.Context, t *testing.T) (p *pgxpool.Pool, teardown func()) {
	helperPool, err := CreateConnectionPool(ctx, testConnectionString(testDBName))
	require.NoError(t, err)
	defer helperPool.Close()

	dbName := getTestDBName(t.Name())
	_, err = helperPool.Exec(ctx, fmt.Sprintf("CREATE DATABASE %s;", dbName))
	require.NoError(t, err)

	pool, err := CreateConnectionPool(ctx, testConnectionString(dbName))
	require.NoError(t, err)
	return pool, func() {
		pool.Close()
	}
}
