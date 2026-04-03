package main

import (
	"database/sql"
	"fmt"
	"github.com/GGnet123/tech_assignment_nanaban/internal/api"
	"github.com/GGnet123/tech_assignment_nanaban/internal/service"
	"github.com/GGnet123/tech_assignment_nanaban/internal/transport/grpc/rates"
	"github.com/GGnet123/tech_assignment_nanaban/pkg/config"
	"github.com/GGnet123/tech_assignment_nanaban/pkg/logger"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	log := logger.New()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Server.Port))
	if err != nil {
		log.Error("listen: %v", err)
		os.Exit(1)
	}

	restyClient := api.NewClient(fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port))
	ratesService := service.NewRateService(restyClient)
	grpcServer := grpc.NewServer()
	rates.Register(grpcServer, log, ratesService)

	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)

	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	if err := grpcServer.Serve(lis); err != nil {
		log.Error("serve: %v", err)
		os.Exit(1)
	}

	// exit channels that listens to process kill
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	<-exit

	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_NOT_SERVING)
	grpcServer.GracefulStop()
}

func checkDB(db *sql.DB, healthServer *health.Server) {
	if err := db.Ping(); err != nil {
		healthServer.SetServingStatus("", healthpb.HealthCheckResponse_NOT_SERVING)
		return
	}

	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
}
