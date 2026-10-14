package http

import (
	"context"
	"github.com/dundduun/msg/core/internal/adapters/http/middleware/auth"
	"github.com/dundduun/msg/core/pkg/response"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
)

type service interface {
	GenTicket(ctx context.Context, uid uuid.UUID) (string, error)
}

type Handler struct {
	log     *slog.Logger
	service service
}

func NewTicketHandler(log *slog.Logger, service service) *Handler {
	return &Handler{
		log:     log,
		service: service,
	}
}

type Response struct {
	Ticket string `json:"ticket,omitempty"`
	response.Response
}

func (h *Handler) GetTicket(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.WithClaims(r.Context())
	if !ok {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, Response{Response: response.Error("can't authorize user")})
		return
	}

	ticket, err := h.service.GenTicket(r.Context(), claims.UID)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, Response{Response: response.Error("failed to gen ticket")})
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, Response{
		Ticket:   ticket,
		Response: response.OK(),
	})
}
