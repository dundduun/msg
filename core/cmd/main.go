package main

import (
	"context"
	"fmt"
	"github.com/dundduun/msg/core/internal/adapters/http/app"
	"github.com/dundduun/msg/core/internal/config"
	"github.com/dundduun/msg/core/internal/profile"
	"github.com/jackc/pgx/v5"
	"log/slog"
	"os"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Name,
	)
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		panic("failed to connect: " + err.Error())
	}

	a := app.New(log, cfg.HTTPServer.Port, conn, profile.NewService())
	a.Start()
}

var (
	envLocal = "local"
	envDev   = "development"
	envProd  = "production"
)

func setupLogger(env string) *slog.Logger {
	var handler slog.Handler

	switch env {
	case envLocal:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	case envDev:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	case envProd:
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	}

	return slog.New(handler)
}
