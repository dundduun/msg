package app

import (
	"context"
	"fmt"
	handlers "github.com/dundduun/msg/core/internal/adapters/http"
	"github.com/dundduun/msg/core/internal/adapters/http/middleware/auth"
	"github.com/dundduun/msg/core/internal/adapters/http/ws"
	"github.com/dundduun/msg/core/internal/infra/postgres"
	cache "github.com/dundduun/msg/core/internal/infra/redis"
	"github.com/dundduun/msg/core/internal/profile"
	"github.com/dundduun/msg/core/internal/ticket"
	"github.com/dundduun/msg/core/pkg/jwt"
	"github.com/dundduun/msg/core/pkg/logerr"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"net/http"
	"time"
)

type App struct {
	log    *slog.Logger
	conn   *pgx.Conn
	rdb    *redis.Client
	server *http.Server
}

func New(log *slog.Logger, rdb *redis.Client, conn *pgx.Conn, port int, profileTTL, ticketTTL time.Duration, secret []byte) *App {
	profileCache, err := cache.NewCache(rdb, profileTTL)
	if err != nil {
		log.Error("failed to create profile cache", logerr.Err(err))
		panic("failed to create profile cache: " + err.Error())
	}
	ticketCache, err := cache.NewTicketCache(rdb, ticketTTL)
	if err != nil {
		log.Error("failed to create ticket cache", logerr.Err(err))
		panic("failed to create ticket cache: " + err.Error())
	}

	profileHandler := handlers.NewProfileHandler(
		profile.NewService(
			profileCache,
			postgres.NewProfileRepo(conn),
			log,
		),
	)
	ticketService := ticket.NewService(ticketCache, log)
	ticketHandler := handlers.NewTicketHandler(log, ticketService)
	wsHandler := ws.NewHandler(log)
	authMiddleware := auth.NewMiddleware(log, jwt.NewService(secret), ticketService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/profile/{id}", profileHandler.GetProfile)
	r.Post("/profile", profileHandler.CreateProfile)
	r.Get("/ws", authMiddleware.EnsureTicket(wsHandler.Connect))
	r.Post("/ws/ticket", authMiddleware.EnsureToken(ticketHandler.GetTicket))

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
	_ = a.rdb.Close()
}
