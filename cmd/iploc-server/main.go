package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/dronnix/search-accomodation/api"
	"github.com/dronnix/search-accomodation/internal/flags"
	"github.com/dronnix/search-accomodation/internal/iploc_api"
	"github.com/dronnix/search-accomodation/storage"
)

type options struct {
	HTTPPort int `long:"http-port" default:"8080" env:"HTTP_PORT"`
	*flags.Postgres
}

func main() {
	opts := &options{}
	flags.Parse(opts)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	storage, err := setupStorage(ctx, opts.Postgres)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not setup storage: %v\n", err)
		return // TODO: set exit code
	}

	ipLocSrv := iploc_api.NewIpLocationServer(storage)
	httpServer := setupHTTPServer(opts, ipLocSrv)

	if err := httpServer.ListenAndServe(); err != nil {
		panic("ListenAndServe: " + err.Error()) // TODO (ALu): replace with logger
	}
	// TODO: Add signal handler
}

func setupStorage(ctx context.Context, opts *flags.Postgres) (*storage.IPLocationStorage, error) {
	pool, err := storage.CreateConnectionPool(ctx, opts.PostgresConnectionString())
	if err != nil {
		return nil, fmt.Errorf("could not create connection pool: %w", err)
	}
	s := storage.NewIPLocationStorage(pool)
	const migrationsPath = "storage/migrations/iplocation" // TODO: move to config
	if err := s.MigrateUp(ctx, migrationsPath); err != nil {
		return nil, fmt.Errorf("could not migrate up: %w", err)
	}
	return s, nil
}

func setupHTTPServer(opts *options, ipLocSrv *iploc_api.IPLocationServer) *http.Server {
	router := chi.NewRouter()
	router.Use(
		middleware.Logger,
		middleware.SetHeader("Content-Type", "application/json"),
		middleware.Heartbeat("/ping"),
		middleware.Recoverer,
	)

	router.Mount("/", api.Handler(ipLocSrv))

	return &http.Server{
		Handler: router,
		Addr:    fmt.Sprintf(":%d", opts.HTTPPort),

		ReadTimeout:    time.Second,
		WriteTimeout:   time.Second,
		IdleTimeout:    time.Minute,
		MaxHeaderBytes: 4096,
	}
}
