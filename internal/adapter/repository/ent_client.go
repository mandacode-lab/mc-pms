package repository

import (
	"errors"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	_ "github.com/lib/pq"
	"github.com/mandacode-com/serengeti-integrated/ent"
)

func NewEntClient(dsn string) (*ent.Client, error) {
	drv, err := sql.Open(dialect.Postgres, dsn)
	if err != nil {
		return nil, err
	}
	if drv == nil {
		return nil, errors.New("failed to create sql driver")
	}

	db := drv.DB()
	if err := db.Ping(); err != nil {
		return nil, err
	}

	client := ent.NewClient(ent.Driver(drv))
	if client == nil {
		return nil, errors.New("failed to create ent client")
	}

	return client, nil
}