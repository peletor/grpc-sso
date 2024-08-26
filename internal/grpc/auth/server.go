package auth

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"google.golang.org/grpc"

	"grpc-sso/internal/grpc/proto/sso"
)

type serverAPI struct {
	sso.UnimplementedAuthServer
	auth Auth
}

type Auth interface {
	Login(ctx context.Context,
		email string,
		password string,
		appID int,
	) (token string, err error)

	RegisterNewUser(ctx context.Context,
		email string,
		password string,
	) (userID int64, err error)

	IsAdmin(ctx context.Context, userID int64) (isAdmin bool, err error)
}

func Register(gRPCServer *grpc.Server, auth Auth) {
	sso.RegisterAuthServer(gRPCServer, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(
	ctx context.Context,
	req *sso.LoginRequest,
) (*sso.LoginResponse, error) {
	if err := validateLogin(req); err != nil {
		return nil, err
	}

	token, err := s.auth.Login(ctx,
		req.GetEmail(),
		req.GetPassword(),
		int(req.GetAppId()))

	if err != nil {
		return nil, status.Error(codes.Internal, "iternal error")
	}

	return &sso.LoginResponse{
		Token: token,
	}, nil
}

func (s *serverAPI) Register(
	ctx context.Context,
	req *sso.RegisterRequest,
) (*sso.RegisterResponse, error) {
	if err := validateRegister(req); err != nil {
		return nil, err
	}

	userID, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, status.Error(codes.Internal, "iternal error")
	}

	return &sso.RegisterResponse{
		UserId: userID,
	}, nil
}

func (s *serverAPI) IsAdmin(
	ctx context.Context,
	req *sso.IsAdminRequest,
) (*sso.IsAdminResponse, error) {
	if err := validateIsAdmin(req); err != nil {
		return nil, err
	}

	isAdmin, err := s.auth.IsAdmin(ctx, req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.Internal, "iternal error")
	}

	return &sso.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil
}

const (
	emptyAppId  = 0
	emptyUserId = 0
)

func validateLogin(req *sso.LoginRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	if req.GetAppId() == emptyAppId {
		return status.Error(codes.InvalidArgument, "appId is required")
	}

	return nil
}

func validateRegister(req *sso.RegisterRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	return nil
}

func validateIsAdmin(req *sso.IsAdminRequest) error {
	if req.GetUserId() == emptyUserId {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	return nil
}
