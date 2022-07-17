package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/dronnix/search-accomodation/api"
	"github.com/dronnix/search-accomodation/internal/flags"
	"github.com/dronnix/search-accomodation/internal/iplocation_api"
	"github.com/dronnix/search-accomodation/storage"
)

type options struct {
	HTTPPort int `long:"http-port" default:"8080" env:"HTTP_PORT"`
	*flags.Postgres
}

const exitCodeOK = 0
const exitCodeError = 1

func main() {
	os.Exit(_main())
}

func _main() int { // separate function to avoid "defer" in main
	opts := &options{}
	flags.Parse(opts)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	storage, err := setupStorage(ctx, opts.Postgres)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not setup storage: %v\n", err)
		return exitCodeError
	}

	ipLocSrv := iplocation_api.NewIpLocationServer(storage)
	httpServer := setupHTTPServer(opts, ipLocSrv)

	setupSignalHandler(ctx, cancel, httpServer) // Gracefully shutdown on SIGINT/SIGTERM.

	if err = httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		return exitCodeError
	}

	fmt.Fprintln(os.Stdout, "server stopped")
	return exitCodeOK
}

// setupStorage connects to the database and performs any necessary migrations.
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

// setupHTTPServer creates and configures the HTTP server and router.
func setupHTTPServer(opts *options, ipLocSrv *iplocation_api.IPLocationServer) *http.Server {
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

func setupSignalHandler(ctx context.Context, cancel func(), apiSrv *http.Server) {
	quitChan := make(chan os.Signal, 1)
	signal.Ignore(syscall.SIGHUP, syscall.SIGPIPE)
	signal.Notify(quitChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quitChan
		if err := apiSrv.Shutdown(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "unable to gracefully shutdown api server")
		}
		cancel()
	}()
}
