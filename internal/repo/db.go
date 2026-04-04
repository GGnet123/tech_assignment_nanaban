package repo

import (
	"context"
	"github.com/GGnet123/tech_assignment_nanaban/internal/domain/rate"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Repository interface {
	BeginTx(ctx context.Context, opts ...pgx.TxOptions) (pgx.Tx, error)
	CommitTx(ctx context.Context, tx pgx.Tx) error
	RollbackTx(ctx context.Context, tx pgx.Tx) error
	SaveRate(ctx context.Context, tx pgx.Tx, r rate.SaveRate) error
}

type DB struct {
	conn *pgxpool.Pool
}

func NewDB(ctx context.Context, dsn string) (*DB, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	// this configs can be moved to .env level
	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = 30 * time.Minute
	config.MaxConnIdleTime = 5 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return &DB{conn: pool}, nil
}

func (db *DB) BeginTx(ctx context.Context, opts ...pgx.TxOptions) (pgx.Tx, error) {
	txOptions := pgx.TxOptions{}
	if len(opts) > 0 {
		txOptions = opts[0]
	}

	return db.conn.BeginTx(ctx, txOptions)
}

func (db *DB) CommitTx(ctx context.Context, tx pgx.Tx) error {
	if tx == nil {
		return nil
	}
	return tx.Commit(ctx)
}

func (db *DB) RollbackTx(ctx context.Context, tx pgx.Tx) error {
	if tx == nil {
		return nil
	}
	return tx.Rollback(ctx)
}

func (db *DB) Close() {
	db.conn.Close()
}
