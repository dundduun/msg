package ws

import (
	"github.com/dundduun/msg/core/pkg/logerr"
	"github.com/gorilla/websocket"
	"log/slog"
	"net/http"
)

type Handler struct {
	log      *slog.Logger
	upgrader websocket.Upgrader
	hub      *Hub
}

func NewHandler(log *slog.Logger) *Handler {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	} // handle check origin for production

	return &Handler{
		log:      log,
		upgrader: upgrader,
		hub:      newHub(),
	}
}

func (h *Handler) Connect(w http.ResponseWriter, r *http.Request) {
	const op = "http.ws.Handler.Connect"
	log := h.log.With(slog.String("op", op)) // add user info

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("failed to upgrade", logerr.Err(err))
		return
	}
	log.Info("connection open")
	defer conn.Close()

	client := &Client{
		conn: conn,
		send: make(chan Envelope, 256),
		hub:  h.hub,
		log:  h.log,
	}
	h.hub.register <- client

	go client.writePump()
	client.readPump()
}
