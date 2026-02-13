package auth

import (
	"context"
	"errors"
	ssov1 "github.com/dundduun/msg/protos/gen/go/sso"
	"github.com/dundduun/msg/sso/internal/services/auth"
	"github.com/dundduun/msg/sso/internal/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(ctx context.Context, email, password string) (token string, err error)
	RegisterNewUser(ctx context.Context, email, password string) (uid int64, err error)
}

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{})
}

//func Register(gRPC *grpc.Server, auth Auth) {
//	ssov1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
//}

func (s *serverAPI) Login(
	ctx context.Context,
	req *ssov1.CredentialsRequest,
) (*ssov1.LoginResponse, error) {
	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.Unauthenticated, "invalid credentials")
		}
		return nil, status.Error(codes.Internal, "failed to login")
	}

	return &ssov1.LoginResponse{Token: token}, nil
}

func (s *serverAPI) Register(
	ctx context.Context,
	req *ssov1.CredentialsRequest,
) (*ssov1.RegisterResponse, error) {
	id, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, storage.ErrEmailTaken) {
			return nil, status.Error(codes.AlreadyExists, "email is already taken")
		}
		return nil, status.Error(codes.Internal, "failed to register")
	}

	return &ssov1.RegisterResponse{
		UserId: id,
	}, nil
}
