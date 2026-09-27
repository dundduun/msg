package http

import (
	"github.com/dundduun/msg/core/pkg/logerr"
	"github.com/gorilla/websocket"
	"log/slog"
	"net/http"
	"time"
)

const (
	pongWait   = 60 * time.Second
	pingPeriod = pongWait * 9 / 10
)

func heartbeat(conn *websocket.Conn) {
	_ = conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		return conn.SetReadDeadline(time.Now().Add(pongWait))
	})
}

type WSHandler struct {
	log      *slog.Logger
	upgrader websocket.Upgrader
}

func NewWSHandler(log *slog.Logger) *WSHandler {
	upgrader := websocket.Upgrader{} // handle check origin for production

	return &WSHandler{
		log:      log,
		upgrader: upgrader,
	}
}

type WSClient struct {
	conn *websocket.Conn
}

func (c *WSClient) ReadPump()  {}
func (c *WSClient) WritePump() {}

func (h *WSHandler) Connect(w http.ResponseWriter, r *http.Request) {
	const op = "http.WSHandler.Connect"
	log := h.log.With(slog.String("op", op)) // add user info

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("failed to upgrade", logerr.Err(err))
		return
	}
	log.Info("connection open")
	defer conn.Close()

	heartbeat(conn)

	go func() {
		ticker := time.NewTicker(pingPeriod)
		defer ticker.Stop()

		for range ticker.C {
			err := conn.WriteMessage(websocket.PingMessage, nil)
			// небезопасно, может попытаться писать одновременно с эхом
			if err != nil {
				log.Error("failed to write ping message", logerr.Err(err))
				return
			}
		}
	}()

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
