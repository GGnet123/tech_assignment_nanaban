package rates

import (
	"github.com/GGnet123/tech_assignment_nanaban/internal/repo"
	"github.com/GGnet123/tech_assignment_nanaban/internal/service"
	"github.com/GGnet123/tech_assignment_nanaban/pkg/logger"
	v1 "github.com/GGnet123/tech_assignment_nanaban/pkg/pb/v1"
	"google.golang.org/grpc"
)

type Server struct {
	v1.UnimplementedRateServiceServer
	rateService *service.Rate
	repo        *repo.DB
	log         *logger.Logger
}

func Register(
	gRPCServer *grpc.Server,
	log *logger.Logger,
	rateService *service.Rate,
	repo *repo.DB,
) {
	v1.RegisterRateServiceServer(gRPCServer, &Server{
		rateService: rateService,
		log:         log,
		repo:        repo,
	})
}
