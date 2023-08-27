package auth

import (
	"github.com/akxcix/passport/pkg/config"
	"github.com/akxcix/passport/pkg/repositories"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"

	_ "github.com/lib/pq"
)

type Database struct {
	db *sqlx.DB
}

func New(conf *config.DatabaseConfig) *Database {
	log.Info().Msg("Connecting to auth database")

	dsn := repositories.FormatPostgresDSN(
		conf.User,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.DatabaseName,
	)
	conn := sqlx.MustConnect("postgres", dsn)
	log.Info().Msg("Connected to auth database")

	return &Database{db: conn}
}
