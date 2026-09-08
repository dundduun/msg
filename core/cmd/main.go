package main

import (
	"context"
	"fmt"
	"github.com/dundduun/msg/core/internal/adapters/http/app"
	"github.com/dundduun/msg/core/internal/config"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	envLocal = "local"
	envDev   = "development"
	envProd  = "production"
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

	log.Info("starting application", slog.Int("port", cfg.HTTPServer.Port), slog.Any("cfg", cfg))
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		panic(err.Error())
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Cache.Host, cfg.Cache.Port),
		Password: cfg.Cache.Password,
		DB:       cfg.Cache.Name,
	})
	if err := rdb.Ping(context.Background()); err != nil {
		panic("failed to set up cache: " + err.Err().Error())
	}

	a := app.New(log, rdb, conn, cfg.HTTPServer.Port)
	go func() {
		a.Start()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	log.Info("stopping")

	ctxTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	a.Stop(ctxTimeout)

	log.Info("app gracefully stopped")
}

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
