package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jessevdk/go-flags"

	"github.com/dronnix/search-accomodation/api"
)

type options struct {
	LogLevel  string `long:"log-level" choice:"trace" choice:"debug" choice:"info" choice:"warn" choice:"error" default:"debug" env:"LOG_LEVEL"` //nolint:lll
	LogFormat string `long:"log-format" choice:"text" choice:"json" default:"text" env:"LOG_FORMAT"`
	HTTPPort  int    `long:"http-port" default:"8080" env:"HTTP_PORT"`
}

func main() {
	opts := &options{}
	parseCfg(opts)

	router := chi.NewRouter()
	router.Use(
		middleware.SetHeader("Content-Type", "application/json"),
		middleware.Heartbeat("/ping"),
		middleware.Recoverer,
	)
	srv := server{}
	router.Mount("/", api.Handler(&srv))

	httpServer := &http.Server{
		Handler: router,
		Addr:    fmt.Sprintf(":%d", opts.HTTPPort),

		ReadTimeout:    time.Second,
		WriteTimeout:   time.Second,
		IdleTimeout:    time.Minute,
		MaxHeaderBytes: 4096,
	}

	if err := httpServer.ListenAndServe(); err != nil {
		panic("ListenAndServe: " + err.Error()) // TODO (ALu): replace with logger
	}
}

func parseCfg(cfg interface{}) {
	parser := flags.NewParser(cfg, flags.Default)
	if _, err := parser.Parse(); err != nil {

		if flagsErr, ok := err.(*flags.Error); ok { //nolint:errorlint
			if flagsErr.Type == flags.ErrHelp {
				os.Exit(0)
			}
			if flagsErr.Type == flags.ErrTag ||
				flagsErr.Type == flags.ErrInvalidTag ||
				flagsErr.Type == flags.ErrDuplicatedFlag ||
				flagsErr.Type == flags.ErrShortNameTooLong {
				panic(err)
			}
		}
		os.Exit(1)
	}
}
