package rates

import (
	"github.com/GGnet123/tech_assignment_nanaban/internal/service"
	v1 "github.com/GGnet123/tech_assignment_nanaban/pkg/pb/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	v1.UnimplementedRateServiceServer
	rateService *service.Rate
	log         *zap.Logger
}

func Register(
	gRPCServer *grpc.Server,
	log *zap.Logger,
	rateService *service.Rate,
) {
	v1.RegisterRateServiceServer(gRPCServer, &Server{
		rateService: rateService,
		log:         log,
	})
}
