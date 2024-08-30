package grpcapp

import (
	"fmt"
	"google.golang.org/grpc"
	grpcauth "grpc-sso/internal/grpc/auth"
	"log/slog"
	"net"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

// New creates new gRPC server app
func New(
	log *slog.Logger,
	authService grpcauth.Auth,
	port int,
) *App {
	gRPCServer := grpc.NewServer()
	grpcauth.Register(gRPCServer, authService)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

// MustRun runs gRPC server and panics if any error occurs
func (app *App) MustRun() {
	err := app.Run()
	if err != nil {
		panic(err)
	}
}

// Run runs gRPC server
func (app *App) Run() error {
	const op = "grpcapp.Run"

	log := app.log.With(slog.String("op", op))

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", app.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("gRPC server is running",
		slog.String("address", listener.Addr().String()),
		slog.Int("port", app.port),
	)

	if err := app.gRPCServer.Serve(listener); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// Stop stops gRPC server
func (app *App) Stop() error {
	const op = "grpcapp.Stop"

	app.log.With(slog.String("op", op)).
		Info("stopping gRPC server", slog.Int("port", app.port))

	app.gRPCServer.GracefulStop()

	return nil
}
