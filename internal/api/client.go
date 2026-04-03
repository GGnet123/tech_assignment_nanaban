package api

import (
	"github.com/go-resty/resty/v2"
	"time"
)

type Resty struct {
	client *resty.Client
}

func NewClient(baseURL string) *Resty {
	client := resty.New().
		SetBaseURL(baseURL).
		SetTimeout(10 * time.Second).
		SetRetryCount(3)

	return &Resty{
		client: client,
	}
}
