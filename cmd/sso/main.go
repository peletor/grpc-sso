package main

import (
	"grpc-sso/internal/app"
	"grpc-sso/internal/config"
	"grpc-sso/internal/logger/slogger"
	"log/slog"
)

func main() {
	cfg := config.MustLoad()

	log := slogger.SetupLogger(cfg.Env)

	log.Info("Starting application", slog.String("config", cfg.Env))

	application := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)

	application.GrpcApp.MustRun()
	// TODO: app init

	// TODO: start gRPC server
}
