package ticket

import (
	"github.com/dundduun/msg/core/internal/adapters/http/middleware/auth"
	"github.com/dundduun/msg/core/pkg/response"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
)

type service interface {
	GenTicket(uid uuid.UUID) (string, error)
}

type Handler struct {
	log     *slog.Logger
	service service
}

type Response struct {
	Ticket string `json:"ticket,omitempty"`
	response.Response
}

func (h *Handler) GetTicket(w http.ResponseWriter, r *http.Request) {
	//const op = "http.ws.ticket.Handler.GetTicket" // to service
	//log := h.log.With(slog.String("op", op)) // to service

	claims, ok := auth.WithClaims(r.Context())
	if !ok {
		//log.Error("can't get claims from context") // to service
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, Response{Response: response.Error("can't authorize user")})
		return
	}
	//log = log.With(slog.Any("uid", claims.UID)) // to service

	ticket, err := h.service.GenTicket(claims.UID)
	if err != nil {
		//log.Error("failed to gen ticket", logerr.Err(err)) // to service
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
