package auth

import (
	//"buf.build/go/protovalidate"
	"context"
	ssov1 "github.com/dundduun/msg/protos/gen/go/sso"
	"google.golang.org/grpc"
	//"google.golang.org/grpc/codes"
	//"google.golang.org/grpc/status"
)

//type Auth interface {
//	Login(email, password string) (token string, err error)
//	Register(email, password string) (err error)
//}

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	//auth Auth
}

func Register(gRPC *grpc.Server) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{})
}

func (s *serverAPI) Login(
	ctx context.Context,
	req *ssov1.CredentialsRequest,
) (*ssov1.LoginResponse, error) {
	return &ssov1.LoginResponse{Token: req.GetEmail() + " " + req.GetPassword()}, nil
	//
	//token, err := s.auth.Login(req.GetEmail(), req.GetPassword())
	//if err != nil {
	//	return nil, status.Errorf(codes.Internal, "failed to login")
	//}
	//
	//return &ssov1.LoginResponse{Token: token}, nil
}

func (s *serverAPI) Register(
	context.Context,
	*ssov1.CredentialsRequest,
) (*ssov1.RegisterResponse, error) {
	panic("implement me")
}
