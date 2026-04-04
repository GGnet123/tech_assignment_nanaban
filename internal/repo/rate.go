package repo

import (
	"context"
	"github.com/GGnet123/tech_assignment_nanaban/internal/domain/rate"
)

func (db *DB) SaveRate(ctx context.Context, rate rate.SaveRate) error {
	_, err := db.conn.Exec(
		ctx,
		"INSERT INTO rates (price, side, timestamp) VALUES ($1, $2, $3)",
		rate.Price, rate.Side, rate.Timestamp,
	)
	if err != nil {
		return err
	}

	return nil
}
