package database

import (
	"errors"
	"log/slog"
	"time"

	"github.com/jmoiron/sqlx"
)

type Options struct {
	Driver string
	DSN    string

	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

var ErrMissingDSN = errors.New("database: missing DSN")

func New(logger *slog.Logger, opts Options) *SDatabase {
	return &SDatabase{
		logger: logger,
		opts:   opts,
	}
}

func NewWithDSN(logger *slog.Logger, dsn string, driver string) *SDatabase {
	switch driver {
	case "mysql":
		driver = "mysql"
	case "postgres":
		driver = "postgres"
	default:
		panic(errors.New("invalid driver"))
	}
	return New(logger, Options{
		Driver: driver,
		DSN:    dsn,
	})
}

func (d *SDatabase) DB() *sqlx.DB {
	return d.db
}
