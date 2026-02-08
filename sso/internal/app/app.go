package app

import (
	grpcapp "github.com/dundduun/msg/sso/internal/app/grpc"
	"go.uber.org/zap"
	"time"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(
	logger *zap.Logger,
	grpcPort int,
	tokenTTL time.Duration,
) *App {
	// TODO: инициализировать storage

	// TODO: инициализировать auth service

	grpcApp := grpcapp.New(logger, grpcPort)

	return &App{
		GRPCSrv: grpcApp,
	}
}

func (a *App) MustRun() {
	a.GRPCSrv.MustRun()
}
