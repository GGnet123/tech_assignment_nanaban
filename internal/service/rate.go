package service

import (
	"context"
	"errors"
	"github.com/GGnet123/tech_assignment_nanaban/internal/api"
	api2 "github.com/GGnet123/tech_assignment_nanaban/internal/domain/api"
	"github.com/GGnet123/tech_assignment_nanaban/internal/domain/rate"
	"github.com/GGnet123/tech_assignment_nanaban/internal/repo"
	v1 "github.com/GGnet123/tech_assignment_nanaban/pkg/pb/v1"
	"golang.org/x/sync/errgroup"
	"strconv"
)

var InvalidNErr = errors.New("invalid N value")
var InvalidMErr = errors.New("invalid M value")

type Rate struct {
	api *api.Resty
	db  *repo.DB
}

func NewRateService(api *api.Resty, db *repo.DB) *Rate {
	return &Rate{api: api, db: db}
}

// Calculate - calculates bid and ask prices, accepts topN and avgNM methods
func (r Rate) Calculate(ctx context.Context, method v1.RateCalcMethod, N, M int) (rate.Result, error) {
	data, err := r.api.GrinexRequest(ctx)
	if err != nil {
		return rate.Result{}, err
	}

	var response rate.Result

	switch method {
	case v1.RateCalcMethod_RATE_CALC_METHOD_TOP_N:
		response, err = getTopN(ctx, N, data)
		if err != nil {
			return rate.Result{}, err
		}
	case v1.RateCalcMethod_RATE_CALC_METHOD_AVG_NM:
		response, err = getAvgNM(ctx, N, M, data)
		if err != nil {
			return rate.Result{}, err
		}
	default:
		// Get top element by default
		response, err = getTopN(ctx, 0, data)
		if err != nil {
			return rate.Result{}, err
		}
	}

	return response, nil
}

func getAvgNM(ctx context.Context, N, M int, data *api2.OrderBook) (rate.Result, error) {
	g, _ := errgroup.WithContext(ctx)
	var response rate.Result
	g.Go(func() error {
		if N >= len(data.Asks) {
			return InvalidNErr
		}

		if M >= len(data.Asks) {
			return InvalidMErr
		}

		total := 0.0
		// including M index
		for _, bid := range data.Bids[N : M+1] {
			val, err := strconv.ParseFloat(bid.Price, 64)
			if err != nil {
				return err
			}
			total += val
		}

		response.Bid = total / float64(M-N)
		return nil
	})

	g.Go(func() error {
		if N >= len(data.Asks) {
			return InvalidNErr
		}

		if M >= len(data.Asks) {
			return InvalidMErr
		}

		total := 0.0
		// including M index
		for _, ask := range data.Asks[N : M+1] {
			val, err := strconv.ParseFloat(ask.Price, 64)
			if err != nil {
				return err
			}
			total += val
		}

		response.Ask = total / float64(M-N)
		return nil
	})

	if err := g.Wait(); err != nil {
		return rate.Result{}, err
	}
	return response, nil
}

func getTopN(ctx context.Context, N int, data *api2.OrderBook) (rate.Result, error) {
	g, _ := errgroup.WithContext(ctx)
	var response rate.Result
	g.Go(func() error {
		if N >= len(data.Bids) {
			return InvalidNErr
		}
		val, err := strconv.ParseFloat(data.Bids[N].Price, 64)
		if err != nil {
			return err
		}
		response.Bid = val
		return nil
	})

	g.Go(func() error {
		if N >= len(data.Asks) {
			return InvalidNErr
		}
		val, err := strconv.ParseFloat(data.Asks[N].Price, 64)
		if err != nil {
			return err
		}
		response.Ask = val
		return nil
	})

	if err := g.Wait(); err != nil {
		return rate.Result{}, err
	}

	return response, nil
}
