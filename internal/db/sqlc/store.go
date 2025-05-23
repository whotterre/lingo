package db

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	Querier
	ExecTx(ctx context.Context, fn func(Querier) error) error
}

// SQLStore holds a connection pool and implements Store
type SQLStore struct {
	pool *pgxpool.Pool
	*Queries
}


func NewSQLStore(pool *pgxpool.Pool) Store {
	return &SQLStore{
		pool:    pool,
		Queries: New(pool),
	}
}

// ExecTx executes a function within a database transaction for ACID compliance.
// It will rollback the transaction if an error occurs
func (store *SQLStore) ExecTx(ctx context.Context, fn func(Querier) error) error {
	tx, err := store.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return rbErr
		}
		return err
	}

	return tx.Commit(ctx)
}
