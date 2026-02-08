package main

import (
	"github.com/dundduun/msg/sso/internal/app"
	"github.com/dundduun/msg/sso/internal/config"
	"go.uber.org/zap"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	logger := mustSetupLogger(cfg.Env)
	defer logger.Sync()

	logger.Info("starting application", zap.Any("cfg", cfg))

	application := app.New(logger, cfg.GRPC.Port, cfg.TokenTTL)
	application.MustRun()

	// TODO: запустить gRPC-сервер
}

func mustSetupLogger(env string) *zap.Logger {
	var log *zap.Logger
	var err error
	switch env {
	case envLocal, envDev:
		log, err = zap.NewDevelopment()
	case envProd:
		log, err = zap.NewProduction()
	}

	if err != nil {
		panic("can't setup logger: " + err.Error())
	}

	return log
}
