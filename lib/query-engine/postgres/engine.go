package postgres

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/multierr"
)

type QueryEngine interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}

type TransactionManager struct {
	pool *pgxpool.Pool
}

type transactionKey string

const key = transactionKey("transaction")

func New(pool *pgxpool.Pool) *TransactionManager {
	return &TransactionManager{pool}
}

func (tm *TransactionManager) GetQueryEngine(ctx context.Context) QueryEngine {
	tx, ok := ctx.Value(key).(QueryEngine)
	if ok && tx != nil {
		return tx
	}

	return tm.pool
}

func (tm *TransactionManager) RunRepeatableRead(ctx context.Context, fx func(ctxTX context.Context) error) error {
	return tm.runWithIsolationLevel(ctx, pgx.RepeatableRead, fx)
}

func (tm *TransactionManager) RunReadCommitted(ctx context.Context, fx func(ctxTX context.Context) error) error {
	return tm.runWithIsolationLevel(ctx, pgx.ReadCommitted, fx)
}

func (tm *TransactionManager) runWithIsolationLevel(ctx context.Context, isoLevel pgx.TxIsoLevel, fx func(ctxTX context.Context) error) error {
	tx, err := tm.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: isoLevel,
	})
	if err != nil {
		return err
	}

	err = fx(context.WithValue(ctx, key, tx))
	if err != nil {
		rollbackErr := tx.Rollback(ctx)
		return multierr.Combine(err, rollbackErr)
	}

	err = tx.Commit(ctx)
	if err != nil {
		rollbackErr := tx.Rollback(ctx)
		return multierr.Combine(err, rollbackErr)
	}

	return nil
}
