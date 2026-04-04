package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/GGnet123/tech_assignment_nanaban/internal/api"
	"github.com/GGnet123/tech_assignment_nanaban/internal/repo"
	"github.com/GGnet123/tech_assignment_nanaban/internal/service"
	"github.com/GGnet123/tech_assignment_nanaban/internal/transport/grpc/rates"
	"github.com/GGnet123/tech_assignment_nanaban/pkg/config"
	"github.com/GGnet123/tech_assignment_nanaban/pkg/logger"
	"github.com/GGnet123/tech_assignment_nanaban/pkg/tracer"
	grpcprom "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	ctx := context.Background()

	log, err := logger.New(cfg.AppEnv)
	defer func() { _ = log.Sync() }()

	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		os.Exit(1)
	}

	shutdown, err := tracer.New(cfg.AppName)
	if err != nil {
		log.Error("Failed to initialize tracer", zap.Error(err))
		os.Exit(1)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Server.Port))
	if err != nil {
		log.Error("Failed to listen", zap.Error(err))
		os.Exit(1)
	}

	db, err := repo.NewDB(ctx, cfg.GetDSN())
	if err != nil {
		log.Error("Failed to create db connection", zap.Error(err))
		os.Exit(1)
	}

	restyClient := api.NewClient()
	ratesService := service.NewRateService(restyClient, db)

	// grpc server with Prometheus
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpcprom.UnaryServerInterceptor),
		grpc.StreamInterceptor(grpcprom.StreamServerInterceptor))

	// register our service
	rates.Register(grpcServer, log, ratesService)

	// init and register health for our server
	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)

	// serving by default
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	reflection.Register(grpcServer)

	// register prometheus
	grpcprom.Register(grpcServer)

	// expose /metrics to access monitoring via http
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	metricsSrv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Prometheus.Port),
		Handler: mux,
	}

	go func() {
		if err := metricsSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("metrics server failed", zap.Error(err))
		}
	}()

	serverErrChan := make(chan error)
	go func() {
		serverErrChan <- grpcServer.Serve(lis)
	}()

	// exit channel that listens to process exit
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrChan:
		log.Error("Server failed", zap.Error(err))
	case <-exit:
		log.Info("shutdown signal received")
	}
	// graceful shutdown
	log.Info("stopping grpc server")
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_NOT_SERVING)
	grpcServer.GracefulStop()

	log.Info("closing db connection")
	db.Close()

	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 5*time.Second)
	defer shutdownCancel()

	log.Info("closing metrics server")
	if err := metricsSrv.Shutdown(shutdownCtx); err != nil {
		log.Error("metrics server failed to gracefully shutdown", zap.Error(err))
	}

	log.Info("closing tracer server")
	if err := shutdown(shutdownCtx); err != nil {
		log.Error("Tracer shutdown failed", zap.Error(err))
	}
}
