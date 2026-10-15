package auth

import (
	"context"
	"errors"
	ticketp "github.com/dundduun/msg/core/internal/ticket"
	"github.com/dundduun/msg/core/pkg/jwt"
	"github.com/dundduun/msg/core/pkg/logerr"
	"github.com/dundduun/msg/core/pkg/response"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"strings"
)

const prefix = "Bearer "

type ctxKey int

const claimsKey ctxKey = iota

func WithClaims(ctx context.Context) (jwt.Claims, bool) {
	claims, ok := ctx.Value(claimsKey).(jwt.Claims)
	return claims, ok
}

type Middleware struct {
	log           *slog.Logger
	tokenService  jwt.Service // to sweep concrete implementation into TokenService.Verify
	ticketService TicketService
}

type TicketService interface {
	Check(ctx context.Context, ticket string) (uuid.UUID, error)
}

func NewMiddleware(log *slog.Logger, tokenService jwt.Service, ticketService TicketService) *Middleware {
	return &Middleware{log: log, tokenService: tokenService, ticketService: ticketService}
}

func (m *Middleware) EnsureToken(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "auth.Middleware.EnsureToken"
		log := m.log.With(slog.String("op", op))

		bearer := r.Header.Get("Authorization")
		tokenStr, found := strings.CutPrefix(bearer, prefix)
		if !found || tokenStr == "" {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, response.Error("bearer token expected"))
			return
		}

		claims, err := m.tokenService.VerifyToken(tokenStr)
		if err != nil {
			problem := "invalid token"
			log.Debug("failed to authorize token", slog.String("problem", problem))
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, response.Error(problem))
			return
		}

		ctx := context.WithValue(r.Context(), claimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func (m *Middleware) EnsureTicket(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "auth.Middleware.EnsureTicket"
		log := m.log.With(slog.String("op", op))

		ticket := r.URL.Query().Get("ticket")
		if ticket == "" {
			problem := "ticket is expected"
			log.Debug(problem)
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, response.Error(problem))
			return
		}

		uid, err := m.ticketService.Check(r.Context(), ticket)
		if err != nil {
			if errors.Is(err, ticketp.ErrTicketNotFound) {
				log.Warn("ticket not found", slog.String("ticket", ticket))
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, response.Error(ticketp.ErrTicketNotFound.Error()))
				return
			}
			log.Error("failed to check ticket", logerr.Err(err))
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, response.Error("failed to check ticket"))
			return
		}
		claims := jwt.Claims{UID: uid}

		ctx := context.WithValue(r.Context(), claimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
