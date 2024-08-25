package main

import (
	"grpc-sso/internal/config"
	"grpc-sso/internal/logger/slogger"
)

func main() {
	cfg := config.MustLoad()

	_ = slogger.SetupLogger(cfg.Env)

	// TODO: app init

	// TODO: start gRPC server
}
