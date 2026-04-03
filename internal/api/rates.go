package api

import (
	"context"
	"fmt"
	"github.com/GGnet123/tech_assignment_nanaban/internal/domain/api"
)

const grinexAPI = "https://grinex.io/api/v1/spot/depth?symbol=usdta7a5"

func (r *Resty) GrinexRequest(ctx context.Context) (*api.OrderBook, error) {
	var result api.OrderBook

	resp, err := r.client.R().
		SetContext(ctx).
		SetResult(&result).
		Get(grinexAPI)

	if err != nil {
		return nil, fmt.Errorf("error while requesting Grinex API: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("grinex error: %s", resp.Status())
	}

	return &result, nil
}
