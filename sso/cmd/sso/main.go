package main

import (
	"github.com/dundduun/msg/sso/internal/app"
	"github.com/dundduun/msg/sso/internal/config"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	logger := mustSetupLogger(cfg.Env)
	defer func() { _ = logger.Sync() }()

	logger.Info("starting application", zap.Any("cfg", cfg))
	application := app.New(logger, cfg.GRPC.Port, cfg.TokenTTL)
	go application.GRPCSrv.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sig := <-stop
	logger.Info("stopping application", zap.String("signal", sig.String()))
	application.GRPCSrv.Stop()

	logger.Info("application stopped")
}

func mustSetupLogger(env string) *zap.Logger {
	var log *zap.Logger
	switch env {
	case envLocal, envDev:
		log = zap.Must(zap.NewDevelopment())
	case envProd:
		log = zap.Must(zap.NewProduction())
	}

	return log
}
