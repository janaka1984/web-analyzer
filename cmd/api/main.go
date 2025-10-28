package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/janaka/web-analyzer/internal/config"
	"github.com/janaka/web-analyzer/internal/transport/httpgin"
	"github.com/janaka/web-analyzer/pkg/logger"
)

func main() {
	cfg := config.MustLoad()
	log := logger.New(cfg.AppEnv)

	r := httpgin.BuildRouter(cfg, log)

	srv := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: r,
	}

	go func() {
		log.Infof("Web Analyzer listening on :%s", cfg.HTTPPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Info("Shutting down gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
}
