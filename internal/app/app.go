package app

import (
	grpcapp "grpc-sso/internal/app/grpc"
	"grpc-sso/internal/services/auth"
	"grpc-sso/internal/storage/sqlite"
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
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, storage, storage, storage, tokenTTL)

	grpcApp := grpcapp.New(log, authService, grpcPort)

	return &App{
		GrpcApp: grpcApp,
	}
}
