package main

import (
	"github.com/GGnet123/tech_assignment_nanaban/pkg/config"
	"github.com/GGnet123/tech_assignment_nanaban/pkg/logger"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func setupRouter(
	cfg *config.Config,
	log *logger.Logger,
) http.Handler {
	r := chi.NewRouter()

	// Health check endpoints
	//healthHandler := health.NewHealthHandler(rateUpdater)
	//r.Get("/health", healthHandler.Health)
	//r.Get("/health/detailed", healthHandler.HealthDetailed)
	//r.Get("/health/live", healthHandler.Live)
	//r.Get("/health/ready", healthHandler.Ready)
	//
	r.Mount("/api/v1", apiV1(cfg, log))

	return r
}

func apiV1(
	cfg *config.Config,
	log *logger.Logger,
) chi.Router {
	r := chi.NewRouter()
	log.Debug("api v1", cfg.Server.Host, cfg.Server.Port)
	//r.Use(middleware.Logger(log))
	//r.Use(middleware.CORS(cfg.App.CORSAllowedOrigins))

	// Public authentication endpoints

	return r
}
