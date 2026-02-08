package auth

import (
	"context"
	ssov1 "github.com/dundduun/msg/protos/gen/go/sso"
	"google.golang.org/grpc"
)

type serverAPI struct {
	ssov1.UnimplementedAuthServer
}

func Register(gRPC *grpc.Server) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{})
}

func (s *serverAPI) Login(
	context.Context,
	*ssov1.CredentialsRequest,
) (*ssov1.LoginResponse, error) {
	panic("implement me")
}

func (s *serverAPI) Register(
	context.Context,
	*ssov1.CredentialsRequest,
) (*ssov1.RegisterResponse, error) {
	panic("implement me")
}
