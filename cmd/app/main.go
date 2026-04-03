package main

import (
	"fmt"
	"github.com/GGnet123/tech_assignment_nanaban/pkg/config"
	"github.com/GGnet123/tech_assignment_nanaban/pkg/logger"
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"time"
)

func main() {
	godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	log := logger.New()
	router := setupRouter(cfg, log)

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	serverErrors := make(chan error, 1)
	go func() {
		serverErrors <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		log.Error("Failed to listen and server", "error", err)
		os.Exit(1)
	}
}
