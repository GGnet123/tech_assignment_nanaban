package rates

import (
	"context"
	v1 "github.com/GGnet123/tech_assignment_nanaban/pkg/pb/v1"
	"go.uber.org/zap"
)

// GetRates - Calls Calculate from rate service and returns the result
func (s *Server) GetRates(ctx context.Context, request *v1.GetRatesRequest) (*v1.GetRatesResponse, error) {
	s.log.Info("GetRates Request Triggered")

	s.log.Info("Calculating rates")
	result, err := s.rateService.Calculate(ctx, request.GetMethod(), int(request.GetN()), int(request.GetM()))
	if err != nil {
		s.log.Error("Calculate rates error", zap.Error(err))
		return nil, err
	}

	s.log.Info("Storing rates", zap.Any("result", result))
	err = s.rateService.SaveRates(ctx, result.Bid, result.Ask, result.Timestamp)
	if err != nil {
		s.log.Error("Failed to store rates", zap.Error(err))
		return nil, err
	}
	return &v1.GetRatesResponse{
		Bid: result.Bid,
		Ask: result.Ask,
	}, nil
}
