package app

import (
	"context"
	"fmt"
	profile "github.com/dundduun/msg/core/internal/adapters/http"
	"github.com/dundduun/msg/core/internal/infra/postgres"
	prof "github.com/dundduun/msg/core/internal/profile"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"log/slog"
	"net/http"
)

type App struct {
	log    *slog.Logger
	conn   *pgx.Conn
	server *http.Server
}

func New(log *slog.Logger, conn *pgx.Conn, port int) *App {
	profileHandler := profile.NewProfileHandler(
		prof.NewService(
			postgres.NewProfileRepo(conn),
			log,
		),
	)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/profile/{id}", profileHandler.GetProfile)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: r,
	}

	return &App{
		log:    log,
		conn:   conn,
		server: server,
	}
}

func (a *App) Start() {
	_ = a.server.ListenAndServe()
}

func (a *App) Stop(ctx context.Context) {
	_ = a.server.Shutdown(ctx)
	_ = a.conn.Close(context.Background())
}
