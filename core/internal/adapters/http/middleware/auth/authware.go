package auth

import (
	"context"
	"github.com/dundduun/msg/core/pkg/jwt"
	"github.com/dundduun/msg/core/pkg/response"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

const prefix = "Bearer "

type ctxKey int

const claimsKey ctxKey = iota

func WithClaims(ctx context.Context) (jwt.Claims, bool) {
	claims, ok := ctx.Value(claimsKey).(jwt.Claims)
	return claims, ok
}

type Middleware struct {
	log     *slog.Logger
	service jwt.Service
}

func NewMiddleware(log *slog.Logger, service jwt.Service) *Middleware {
	return &Middleware{log: log, service: service}
}

func (a *Middleware) Ensure(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.URL.Query().Get("ticket")
		if tokenString == "" {
			problem := "ticket is expected"
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, response.Error(problem))
			return
		}

		claims, err := a.service.VerifyToken(tokenString)
		if err != nil {
			problem := "invalid ticket"
			a.log.Debug("failed to authorize ticket", slog.String("problem", problem))
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, response.Error(problem))
			return
		}

		ctx := context.WithValue(r.Context(), claimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
