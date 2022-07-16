package flags

import (
	"fmt"
	"net/url"
	"os"

	"github.com/jessevdk/go-flags"
)

type Postgres struct {
	PostgresHost string `long:"postgres-host" description:"host with PG" default:"localhost" env:"POSTGRES_HOST"`
	PostgresPort string `long:"postgres-port" description:"port where PG is listening" default:"5432" env:"POSTGRES_PORT"`
	PostgresDB   string `long:"postgres-db" description:"name of the database to connect to" default:"geolocation" env:"POSTGRES_DB"` // nolint:lll
	PostgresUser string `long:"postgres-user" description:"PG user" default:"postgres" env:"POSTGRES_USER"`
	PostgresPass string `long:"postgres-pass" description:"PG password" default:"NA" env:"POSTGRES_PASS"`
}

func Parse(cfg interface{}) {
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

func (p *Postgres) PostgresConnectionString() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		url.QueryEscape(p.PostgresUser), url.QueryEscape(p.PostgresPass),
		p.PostgresHost, p.PostgresPort, p.PostgresDB,
	)
}
