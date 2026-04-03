package main

import (
	"context"
	"fmt"
	"github.com/GGnet123/tech_assignment_nanaban/internal/api"
	"github.com/GGnet123/tech_assignment_nanaban/internal/repo"
	"github.com/GGnet123/tech_assignment_nanaban/internal/service"
	"github.com/GGnet123/tech_assignment_nanaban/internal/transport/grpc/rates"
	"github.com/GGnet123/tech_assignment_nanaban/pkg/config"
	"github.com/GGnet123/tech_assignment_nanaban/pkg/logger"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log := logger.New()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Server.Port))
	if err != nil {
		log.Error("listen: %v", err)
		os.Exit(1)
	}

	db, err := repo.NewDB(ctx, cfg.GetDSN())
	if err != nil {
		log.Error("Couldn't create db connection: %v", err)
		os.Exit(1)
	}

	restyClient := api.NewClient()
	ratesService := service.NewRateService(restyClient, db)

	// grpc server
	grpcServer := grpc.NewServer()
	// register our service
	rates.Register(grpcServer, log, ratesService, db)

	// init and register health for our server
	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)

	// serving by default
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	reflection.Register(grpcServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.Error("serve: %v", err)
		os.Exit(1)
	}

	// exit channels that listens to process exit
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	<-exit

	// graceful shutdown
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_NOT_SERVING)
	grpcServer.GracefulStop()

	db.Close()
}
