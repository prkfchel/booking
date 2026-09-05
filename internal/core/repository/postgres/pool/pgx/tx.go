package core_pgx_pool

import (
	core_postgres_pool "booking/internal/core/repository/postgres/pool"
	"context"

	"github.com/jackc/pgx/v5"
)

type pgxTx struct {
	pgx.Tx
}

func (tx *pgxTx) Commit(ctx context.Context) error {
	err := tx.Tx.Commit(ctx)
	if err != nil {
		return mapErrors(err)
	}

	return nil
}

func (tx *pgxTx) Rollback(ctx context.Context) error {
	err := tx.Tx.Rollback(ctx)
	if err != nil {
		return mapErrors(err)
	}

	return nil
}

func (tx *pgxTx) Query(ctx context.Context, sql string, args ...any) (core_postgres_pool.Rows, error) {
	rows, err := tx.Tx.Query(ctx, sql, args...)
	if err != nil {
		return pgxRows{}, mapErrors(err)
	}

	return rows, nil
}

func (tx *pgxTx) QueryRow(ctx context.Context, sql string, args ...any) core_postgres_pool.Row {
	row := tx.Tx.QueryRow(ctx, sql, args...)
	return row
}

func (tx *pgxTx) Exec(ctx context.Context, sql string, arguments ...any) (core_postgres_pool.CommandTag, error) {
	commandTag, err := tx.Tx.Exec(ctx, sql, arguments...)
	if err != nil {
		return pgxCommandTag{}, mapErrors(err)
	}

	return commandTag, nil
}
