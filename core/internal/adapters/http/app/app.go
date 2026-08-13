package app

import (
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
	log            *slog.Logger
	port           int
	profileHandler *profile.ProfileHandler
}

func New(log *slog.Logger, conn *pgx.Conn, port int) *App {
	return &App{
		log:  log,
		port: port,
		profileHandler: profile.NewProfileHandler(
			prof.NewService(
				postgres.NewProfileRepo(conn),
				log,
			),
		),
	}
}

func (a *App) Start() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/profile/{id}", a.profileHandler.GetProfile)

	a.log.Info("starting application")
	_ = http.ListenAndServe(fmt.Sprintf(":%d", a.port), r)
	a.log.Info("server stopped")
}
