package main

import (
	"grpc-sso/internal/app"
	"grpc-sso/internal/config"
	"grpc-sso/internal/logger/slogger"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()

	log := slogger.SetupLogger(cfg.Env)

	log.Info("Starting application", slog.String("config", cfg.Env))

	application := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)

	go application.GrpcApp.MustRun()

	// Graceful shutdown

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	sign := <-stop

	log.Info("Stopping application", slog.String("signal", sign.String()))

	application.GrpcApp.Stop()

	log.Info("Application stopped")
}
