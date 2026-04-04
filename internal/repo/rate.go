package repo

import (
	"context"
	"github.com/GGnet123/tech_assignment_nanaban/internal/domain/rate"
	"github.com/jackc/pgx/v5"
)

func (db *DB) SaveRate(ctx context.Context, tx pgx.Tx, r rate.SaveRate) error {
	_, err := tx.Exec(
		ctx,
		"INSERT INTO rates (price, side, timestamp) VALUES ($1, $2, $3)",
		r.Price, r.Side, r.Timestamp,
	)
	return err
}
