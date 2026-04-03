package rates

import (
	"context"
	v1 "github.com/GGnet123/tech_assignment_nanaban/pkg/pb/v1"
	"go.uber.org/zap"
)

// GetRates - Calls Calculate from rate service and returns the result
func (s *Server) GetRates(ctx context.Context, request *v1.GetRatesRequest) (*v1.GetRatesResponse, error) {
	result, err := s.rateService.Calculate(ctx, request.GetMethod(), int(request.GetN()), int(request.GetM()))
	if err != nil {
		s.log.Error("Calculate rates error", zap.Error(err))
		return nil, err
	}
	return &v1.GetRatesResponse{
		Bid: result.Bid,
		Ask: result.Ask,
	}, nil
}
