package auth

import (
	"context"

	"google.golang.org/grpc"

	"grpc-sso/internal/grpc/proto/sso"
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

func Register(gRPCServer *grpc.Server) {
	sso.RegisterAuthServer(gRPCServer, &serverAPI{ /*auth: auth*/ })
}

func (s *serverAPI) Login(
	ctx context.Context,
	in *sso.LoginRequest,
) (*sso.LoginResponse, error) {
	panic("implement me")
}

func (s *serverAPI) Register(
	ctx context.Context,
	req *sso.RegisterRequest,
) (*sso.RegisterResponse, error) {
	panic("implement me")
}

func (s *serverAPI) IsAdmin(
	ctx context.Context,
	req *sso.IsAdminRequest,
) (*sso.IsAdminResponse, error) {
	panic("implement me")
}
