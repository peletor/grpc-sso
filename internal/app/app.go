package app

import (
	grpcapp "grpc-sso/internal/app/grpc"
	"log/slog"
	"time"
)

type App struct {
	GrpcApp *grpcapp.App
}

// New creates new gRPC server app
func New(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
	tokenTTL time.Duration,
) *App {
	// TODO: init storage

	// TODO: init auth service
	grpcApp := grpcapp.New(log, grpcPort)

	return &App{
		GrpcApp: grpcApp,
	}
}
