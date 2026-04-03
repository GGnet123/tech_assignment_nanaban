package rates

import (
	"context"
	"github.com/GGnet123/tech_assignment_nanaban/internal/domain/rate"
	v1 "github.com/GGnet123/tech_assignment_nanaban/pkg/pb/v1"
)

// GetRates - Calls Calculate from rate service and returns the result
func (s *Server) GetRates(ctx context.Context, request *v1.GetRatesRequest) (*v1.GetRatesResponse, error) {
	s.log.Info("GetRates Request Triggered")
	s.log.Info("Calculation rates")
	result, err := s.rateService.Calculate(ctx, request.GetMethod(), int(request.GetN()), int(request.GetM()))
	if err != nil {
		s.log.Error("Calculate rates error", err)
		return nil, err
	}

	s.log.Info("Storing rates", result)
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		s.log.Error("Create tx error", err)
		return nil, err
	}

	err = s.repo.SaveRate(ctx, rate.SaveRate{
		Price: result.Bid,
		Side:  rate.SideBid,
	})

	if err != nil {
		s.log.Error("Save rates bid error", err)
		return nil, err
	}

	err = s.repo.SaveRate(ctx, rate.SaveRate{
		Price: result.Ask,
		Side:  rate.SideAsk,
	})

	if err != nil {
		s.log.Error("Save rates ask error", err)
		return nil, err
	}

	err = s.repo.CommitTx(ctx, tx)
	if err != nil {
		s.log.Error("Commit tx error", err)
		return nil, err
	}

	return &v1.GetRatesResponse{
		Bid: result.Bid,
		Ask: result.Ask,
	}, nil
}
