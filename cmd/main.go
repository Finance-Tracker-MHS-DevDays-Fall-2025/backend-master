package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"backend-master/configs"
	"backend-master/internal"
	"backend-master/internal/logger"

	"go.uber.org/zap"
)

func main() {
	cfg, err := configs.New()
	if err != nil {
		panic(err)
	}

	logger, err := logger.NewLogger()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	service := internal.NewService(cfg, logger)

	go func() {
		err := service.Start()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("failed to start service", zap.Error(err))
		}
	}()

	logger.Info("server started successfully")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 10)
	defer cancel()

	if err := service.Shutdown(); err != nil {
		logger.Error("error during shutdown", zap.Error(err))
	}

	<-ctx.Done()
}
