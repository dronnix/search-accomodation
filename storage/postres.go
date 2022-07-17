package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

func CreateConnectionPool(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("cannot parse connection string: %w", err)
	}

	pool, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("unable to connect: %w", err)
	}
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to accuire connection: %w", err)
	}
	defer conn.Release()
	if err := conn.Conn().Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping DB: %w", err)
	}
	return pool, nil
}
