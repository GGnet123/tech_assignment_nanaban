package api

import (
	"github.com/go-resty/resty/v2"
	"time"
)

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
