package api

import (
	"context"
	"github.com/GGnet123/tech_assignment_nanaban/internal/domain/api"
	"github.com/go-resty/resty/v2"
	"time"
)

type ApiClient interface {
	GrinexRequest(ctx context.Context) (*api.OrderBook, error)
}

type Resty struct {
	client *resty.Client
}

func NewClient() *Resty {
	client := resty.New().
		SetTimeout(10 * time.Second).
		SetRetryCount(3)

	return &Resty{
		client: client,
	}
}
