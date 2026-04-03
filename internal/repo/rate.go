package repo

import (
	"context"
	"github.com/GGnet123/tech_assignment_nanaban/internal/domain/rate"
)

func (db *DB) SaveRate(ctx context.Context, rate rate.SaveRate) error {
	_, err := db.conn.Exec(
		ctx,
		"INSERT INTO rates (price, side) VALUES ($1, $2)",
		rate.Price, rate.Side,
	)
	if err != nil {
		return err
	}

	return nil
}
