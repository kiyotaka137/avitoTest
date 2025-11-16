package main

import (
	"context"
	"time"

	"avitoTest/internal/app"
	"avitoTest/internal/config"
	"avitoTest/pkg/logger"

	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	logger.InitLogger(cfg)
	logger.Log.Info("configuration loaded", "log_level", cfg.Log.Level)

	a := app.New(logger.Log, cfg)

	go func() {
		if err := a.Run(); err != nil {
			logger.Log.Error("http server stopped", "err", err)
		}
	}()

	logger.Log.Info("application started", "addr", cfg.Http.Address)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	logger.Log.Info("shutdown signal received")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := a.Shutdown(ctx); err != nil {
		logger.Log.Error("graceful shutdown failed", "err", err)
	} else {
		logger.Log.Info("server stopped gracefully")
	}
}
