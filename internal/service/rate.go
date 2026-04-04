package service

import (
	"context"
	"errors"
	"github.com/GGnet123/tech_assignment_nanaban/internal/api"
	api2 "github.com/GGnet123/tech_assignment_nanaban/internal/domain/api"
	"github.com/GGnet123/tech_assignment_nanaban/internal/domain/rate"
	"github.com/GGnet123/tech_assignment_nanaban/internal/repo"
	v1 "github.com/GGnet123/tech_assignment_nanaban/pkg/pb/v1"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"golang.org/x/sync/errgroup"
	"math"
	"strconv"
)

var ErrInvalidN = errors.New("invalid N value")
var ErrInvalidM = errors.New("invalid M value")

type Rate struct {
	api  api.ApiClient
	repo repo.Repository
}

func NewRateService(api api.ApiClient, repo repo.Repository) *Rate {
	return &Rate{api: api, repo: repo}
}

func (r Rate) SaveRates(ctx context.Context, bid, ask float64, timestamp int64) error {
	ctx, span := otel.Tracer("rate").Start(ctx, "SaveRates")
	defer span.End()

	// error handling function
	fail := func(err error, tx pgx.Tx) error {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		if rbErr := r.repo.RollbackTx(ctx, tx); rbErr != nil {
			span.RecordError(rbErr)
			return errors.Join(err, rbErr)
		}
		return err
	}

	tx, err := r.repo.BeginTx(ctx)
	if err != nil {
		return fail(err, tx)
	}

	err = r.repo.SaveRate(ctx, tx, rate.SaveRate{
		Price:     bid,
		Side:      rate.SideBid,
		Timestamp: timestamp,
	})

	if err != nil {
		return fail(err, tx)
	}

	err = r.repo.SaveRate(ctx, tx, rate.SaveRate{
		Price:     ask,
		Side:      rate.SideAsk,
		Timestamp: timestamp,
	})

	if err != nil {
		return fail(err, tx)
	}

	err = r.repo.CommitTx(ctx, tx)
	if err != nil {
		return fail(err, tx)
	}

	return nil
}

// Calculate - calculates bid and ask prices, accepts topN and avgNM methods
func (r Rate) Calculate(ctx context.Context, method v1.RateCalcMethod, N, M int) (rate.Result, error) {
	ctx, span := otel.Tracer("rate").Start(ctx, "Calculate")
	defer span.End()

	data, err := r.api.GrinexRequest(ctx)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return rate.Result{}, err
	}

	var response rate.Result
	switch method {
	case v1.RateCalcMethod_RATE_CALC_METHOD_TOP_N:
		response, err = getTopN(ctx, N, data)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return rate.Result{}, err
		}
	case v1.RateCalcMethod_RATE_CALC_METHOD_AVG_NM:
		response, err = getAvgNM(ctx, N, M, data)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return rate.Result{}, err
		}
	default:
		// Get top element by default
		response, err = getTopN(ctx, 0, data)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return rate.Result{}, err
		}
	}

	response.Timestamp = data.Timestamp
	return response, nil
}

func getAvgNM(ctx context.Context, N, M int, data *api2.OrderBook) (rate.Result, error) {
	g, _ := errgroup.WithContext(ctx)
	var response rate.Result
	g.Go(func() error {
		if N >= len(data.Bids) {
			return ErrInvalidN
		}

		if M >= len(data.Bids) {
			return ErrInvalidM
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

		avg := total / float64(M-N+1)

		// round to have only 2 digits after comma
		response.Bid = math.Round(avg*100) / 100
		return nil
	})

	g.Go(func() error {
		if N >= len(data.Asks) {
			return ErrInvalidN
		}

		if M >= len(data.Asks) {
			return ErrInvalidM
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
		avg := total / float64(M-N+1)
		// round to have only 2 digits after comma
		response.Ask = math.Round(avg*100) / 100
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
			return ErrInvalidN
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
			return ErrInvalidN
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
