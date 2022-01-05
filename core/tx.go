package core

import (
	"context"
	"database/sql"
)

type Tx struct {
	*sql.Tx
	db  *DB
	ctx context.Context
}

func (db *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	tx, err := db.DB.BeginTx(ctx, opts)
	return &Tx{tx, db, ctx}, err
}

func (db *DB) Begin() (*Tx, error) {
	return db.BeginTx(context.Background(), nil)
}

func (tx *Tx) Commit() error {
	return tx.Tx.Commit()
}

func (tx *Tx) Rollback() error {
	return tx.Tx.Rollback()
}

func (tx *Tx) Exec(query string, args ...interface{}) (sql.Result, error) {
	return tx.ExecContext(tx.ctx, query, args...)
}

func (tx *Tx) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	res, err := tx.Tx.ExecContext(ctx, query, args...)
	return res, err
}

// Query query with args
func (tx *Tx) Query(query string, args ...interface{}) (*Rows, error) {
	return tx.QueryContext(tx.ctx, query, args...)
}

// QueryContext query with args
func (tx *Tx) QueryContext(ctx context.Context, query string, args ...interface{}) (*Rows, error) {
	rows, err := tx.Tx.QueryContext(ctx, query, args...)
	return &Rows{rows, tx.db}, err
}
