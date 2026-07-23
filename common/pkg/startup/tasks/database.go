package tasks

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	db *pgxpool.Pool
}

func NewDatabase(db *pgxpool.Pool) *Database {

	return &Database{
		db: db,
	}
}

func (d *Database) Name() string {
	return "database"
}

func (d *Database) Run(ctx context.Context) error {
	return d.db.Ping(ctx)
}
