package main

import "grpc-sso/internal/config"

func main() {
	config.MustLoad()

	// TODO: logger init

	// TODO: app init

	// TODO: start gRPC servergit
}
