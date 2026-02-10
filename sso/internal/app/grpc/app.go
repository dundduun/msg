package grpcapp

import (
	"fmt"
	authgrpc "github.com/dundduun/msg/sso/internal/grpc/auth"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

type App struct {
	logger     *zap.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(logger *zap.Logger, port int) *App {
	srv := grpc.NewServer()

	authgrpc.Register(srv)

	return &App{
		logger:     logger,
		gRPCServer: srv,
		port:       port,
	}
}

func (a *App) MustRun() {
	err := a.Run()
	if err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "grpcapp.Run"

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	a.logger.Info("gRPC server is running", zap.String("addr", l.Addr().String()))

	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.logger.With(zap.String("op", op)).
		Info("stopping gRPC server")

	a.gRPCServer.GracefulStop()
}
