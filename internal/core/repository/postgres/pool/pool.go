package core_postgres_pool

import (
	"context"
	"time"
)

type Pool interface {
	Query(ctx context.Context, sql string, args ...any) (Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) Row
	Exec(ctx context.Context, sql string, arguments ...any) (CommandTag, error)
	Close()
	OpTimeout() time.Duration
	Begin(ctx context.Context) (Tx, error)
	BeginTx(ctx context.Context, txOptions TxOptions) (Tx, error)
}

type TxOptions struct {
	IsoLevel   TxIsoLevel
	AccessMode TxAccessMode
}

type TxIsoLevel string
type TxAccessMode string

type Tx interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	Query(ctx context.Context, sql string, args ...any) (Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) Row
	Exec(ctx context.Context, sql string, arguments ...any) (CommandTag, error)
}

type Rows interface {
	Close()
	Err() error
	Next() bool
	Scan(dest ...any) error
}

type Row interface {
	Scan(dest ...any) error
}

type CommandTag interface {
	RowsAffected() int64
}
