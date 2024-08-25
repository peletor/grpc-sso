package auth

import (
	"context"

	"google.golang.org/grpc"

	sso "grpc-sso/protos/gen/go"
)

type serverAPI struct {
	sso.UnimplementedAuthServer
	auth Auth
}

type Auth interface {
	Login(
		ctx context.Context,
		email string,
		password string,
		appID int,
	) (token string, err error)
	RegisterNewUser(
		ctx context.Context,
		email string,
		password string,
	) (userID int64, err error)
}

func Register(gRPCServer *grpc.Server, auth Auth) {
	sso.RegisterAuthServer(gRPCServer, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(
	ctx context.Context,
	in *sso.LoginRequest,
) (*sso.LoginResponse, error) {
	// TODO
}

func (s *serverAPI) Register(
	ctx context.Context,
	in sso.RegisterRequest,
) (*sso.RegisterResponce, err error) {
}
