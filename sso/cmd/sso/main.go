package main

import (
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

	_ = setupLogger(cfg.Env)

	// TODO: запустить приложение

	// TODO: запустить gRPC-сервер
}

func setupLogger(env string) *zap.Logger {
	var log *zap.Logger
	switch env {
	case envLocal, envDev:
		log, _ = zap.NewDevelopment()
	case envProd:
		log, _ = zap.NewProduction()
	}

	return log
}
