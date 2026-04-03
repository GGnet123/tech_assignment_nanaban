package repo

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

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

func (db *DB) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return db.conn.Begin(ctx)
}

func (db *DB) CommitTx(ctx context.Context, tx pgx.Tx) error {
	return tx.Commit(ctx)
}

func (db *DB) Close() {
	db.conn.Close()
}
