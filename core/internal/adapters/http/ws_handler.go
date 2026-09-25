package http

import (
	"github.com/dundduun/msg/core/pkg/logerr"
	"github.com/gorilla/websocket"
	"log/slog"
	"net/http"
)

type WSHandler struct {
	log      *slog.Logger
	upgrader websocket.Upgrader
}

func NewChatHandler(log *slog.Logger) *WSHandler {
	upgrader := websocket.Upgrader{} // handle check origin for production

	return &WSHandler{
		log:      log,
		upgrader: upgrader,
	}
}

func (c *WSHandler) Connect(w http.ResponseWriter, r *http.Request) {
	const op = "http.WSHandler.Connect"
	log := c.log.With(slog.String("op", op)) // add user info

	conn, err := c.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("failed to upgrade", logerr.Err(err))
		return
	}
	log.Info("connection open")

	defer conn.Close()
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Error("failed to read message", logerr.Err(err))
			} else {
				log.Info("connection closed", slog.String("reason", err.Error()))
			}
			break
		}

		err = conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Error("failed to write message", logerr.Err(err))
			break
		}
	}
}
