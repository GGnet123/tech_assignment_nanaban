package service

import (
	"context"
	"errors"
	"testing"

	"github.com/GGnet123/tech_assignment_nanaban/internal/domain/api"
	"github.com/GGnet123/tech_assignment_nanaban/internal/domain/rate"
	v1 "github.com/GGnet123/tech_assignment_nanaban/pkg/pb/v1"
	"github.com/jackc/pgx/v5"
)

// --- mocks ---

type mockAPI struct {
	result *api.OrderBook
	err    error
}

func (m *mockAPI) GrinexRequest(_ context.Context) (*api.OrderBook, error) {
	return m.result, m.err
}

type mockRepo struct {
	beginErr    error
	saveErr     error
	commitErr   error
	rollbackErr error
	savedRates  []rate.SaveRate
}

func (m *mockRepo) BeginTx(_ context.Context, _ ...pgx.TxOptions) (pgx.Tx, error) {
	return nil, m.beginErr
}

func (m *mockRepo) CommitTx(_ context.Context, _ pgx.Tx) error {
	return m.commitErr
}

func (m *mockRepo) RollbackTx(_ context.Context, _ pgx.Tx) error {
	return m.rollbackErr
}

func (m *mockRepo) SaveRate(_ context.Context, _ pgx.Tx, r rate.SaveRate) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	m.savedRates = append(m.savedRates, r)
	return nil
}

// --- helpers ---

func orderBook(bids, asks []string) *api.OrderBook {
	toOrders := func(prices []string) []api.Order {
		orders := make([]api.Order, len(prices))
		for i, p := range prices {
			orders[i] = api.Order{Price: p}
		}
		return orders
	}
	return &api.OrderBook{Bids: toOrders(bids), Asks: toOrders(asks), Timestamp: 1000}
}

// --- Calculate tests ---

func TestCalculate_TopN(t *testing.T) {
	svc := NewRateService(&mockAPI{
		result: orderBook(
			[]string{"100.0", "99.0", "98.0"},
			[]string{"101.0", "102.0", "103.0"},
		),
	}, &mockRepo{})

	result, err := svc.Calculate(context.Background(), v1.RateCalcMethod_RATE_CALC_METHOD_TOP_N, 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Bid != 100.0 {
		t.Errorf("expected bid 100.0, got %v", result.Bid)
	}
	if result.Ask != 101.0 {
		t.Errorf("expected ask 101.0, got %v", result.Ask)
	}
	if result.Timestamp != 1000 {
		t.Errorf("expected timestamp 1000, got %v", result.Timestamp)
	}
}

func TestCalculate_TopN_OutOfBounds(t *testing.T) {
	svc := NewRateService(&mockAPI{
		result: orderBook([]string{"100.0"}, []string{"101.0"}),
	}, &mockRepo{})

	_, err := svc.Calculate(context.Background(), v1.RateCalcMethod_RATE_CALC_METHOD_TOP_N, 5, 0)
	if !errors.Is(err, ErrInvalidN) {
		t.Errorf("expected ErrInvalidN, got %v", err)
	}
}

func TestCalculate_AvgNM(t *testing.T) {
	svc := NewRateService(&mockAPI{
		result: orderBook(
			[]string{"100.0", "200.0", "300.0"},
			[]string{"110.0", "210.0", "310.0"},
		),
	}, &mockRepo{})

	// avg of index 0..2 = (100+200+300)/3 = 200
	result, err := svc.Calculate(context.Background(), v1.RateCalcMethod_RATE_CALC_METHOD_AVG_NM, 0, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Bid != 200.0 {
		t.Errorf("expected bid 200.0, got %v", result.Bid)
	}
	if result.Ask != 210.0 {
		t.Errorf("expected ask 210.0, got %v", result.Ask)
	}
}

func TestCalculate_Default_UsesTopZero(t *testing.T) {
	svc := NewRateService(&mockAPI{
		result: orderBook([]string{"50.0"}, []string{"51.0"}),
	}, &mockRepo{})

	result, err := svc.Calculate(context.Background(), v1.RateCalcMethod_RATE_CALC_METHOD_UNSPECIFIED, 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Bid != 50.0 {
		t.Errorf("expected bid 50.0, got %v", result.Bid)
	}
}

func TestCalculate_APIError(t *testing.T) {
	apiErr := errors.New("api unavailable")
	svc := NewRateService(&mockAPI{err: apiErr}, &mockRepo{})

	_, err := svc.Calculate(context.Background(), v1.RateCalcMethod_RATE_CALC_METHOD_TOP_N, 0, 0)
	if !errors.Is(err, apiErr) {
		t.Errorf("expected api error, got %v", err)
	}
}

// --- SaveRates tests ---

func TestSaveRates_Success(t *testing.T) {
	repo := &mockRepo{}
	svc := NewRateService(&mockAPI{}, repo)

	err := svc.SaveRates(context.Background(), 100.0, 101.0, 1000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(repo.savedRates) != 2 {
		t.Fatalf("expected 2 saved rates, got %d", len(repo.savedRates))
	}
	if repo.savedRates[0].Side != rate.SideBid || repo.savedRates[0].Price != 100.0 {
		t.Errorf("unexpected bid rate: %+v", repo.savedRates[0])
	}
	if repo.savedRates[1].Side != rate.SideAsk || repo.savedRates[1].Price != 101.0 {
		t.Errorf("unexpected ask rate: %+v", repo.savedRates[1])
	}
}

func TestSaveRates_BeginTxError(t *testing.T) {
	txErr := errors.New("begin failed")
	svc := NewRateService(&mockAPI{}, &mockRepo{beginErr: txErr})

	err := svc.SaveRates(context.Background(), 100.0, 101.0, 1000)
	if !errors.Is(err, txErr) {
		t.Errorf("expected begin tx error, got %v", err)
	}
}

func TestSaveRates_SaveError_Rollback(t *testing.T) {
	saveErr := errors.New("save failed")
	svc := NewRateService(&mockAPI{}, &mockRepo{saveErr: saveErr})

	err := svc.SaveRates(context.Background(), 100.0, 101.0, 1000)
	if !errors.Is(err, saveErr) {
		t.Errorf("expected save error, got %v", err)
	}
}

func TestSaveRates_SaveError_RollbackFails(t *testing.T) {
	saveErr := errors.New("save failed")
	rbErr := errors.New("rollback failed")
	svc := NewRateService(&mockAPI{}, &mockRepo{saveErr: saveErr, rollbackErr: rbErr})

	err := svc.SaveRates(context.Background(), 100.0, 101.0, 1000)
	if !errors.Is(err, saveErr) {
		t.Errorf("expected save error in joined error, got %v", err)
	}
	if !errors.Is(err, rbErr) {
		t.Errorf("expected rollback error in joined error, got %v", err)
	}
}
