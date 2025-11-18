package repo

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type DBExecutor interface {
	WithTx(ctx context.Context, fn func(tx *sqlx.Tx) error) error
}

type dbExecutor struct {
	db *sqlx.DB
}

func NewDBExecutor(db *sqlx.DB) DBExecutor {
	return &dbExecutor{db: db}
}

func (d *dbExecutor) WithTx(ctx context.Context, fn func(tx *sqlx.Tx) error) error {
	tx, err := d.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
