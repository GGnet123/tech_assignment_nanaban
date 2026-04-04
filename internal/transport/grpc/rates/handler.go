package rates

import (
	"context"
	v1 "github.com/GGnet123/tech_assignment_nanaban/pkg/pb/v1"
	"log/slog"
)

// GetRates - Calls Calculate from rate service and returns the result
func (s *Server) GetRates(ctx context.Context, request *v1.GetRatesRequest) (*v1.GetRatesResponse, error) {
	s.log.Info("GetRates Request Triggered")
	s.log.Info("Calculation rates")
	result, err := s.rateService.Calculate(ctx, request.GetMethod(), int(request.GetN()), int(request.GetM()))
	if err != nil {
		s.log.Error("Calculate rates error", slog.Any("error", err))
		return nil, err
	}

	s.log.Info("Storing rates", slog.Any("result", result))
	err = s.rateService.SaveRates(ctx, result.Bid, result.Ask, result.Timestamp)
	if err != nil {
		s.log.Info("Failed to store rates")
		return nil, err
	}
	return &v1.GetRatesResponse{
		Bid: result.Bid,
		Ask: result.Ask,
	}, nil
}
