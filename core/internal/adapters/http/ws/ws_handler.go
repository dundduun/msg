package ws

import (
	"github.com/dundduun/msg/core/pkg/logerr"
	"github.com/gorilla/websocket"
	"log/slog"
	"net/http"
	"time"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = pongWait * 9 / 10
)

func heartbeat(conn *websocket.Conn) {
	_ = conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		return conn.SetReadDeadline(time.Now().Add(pongWait))
	})
}

type Handler struct {
	log      *slog.Logger
	upgrader websocket.Upgrader
	hub      *Hub
}

func NewHandler(log *slog.Logger) *Handler {
	upgrader := websocket.Upgrader{} // handle check origin for production

	return &Handler{
		log:      log,
		upgrader: upgrader,
		hub:      newHub(),
	}
}

type Hub struct {
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	clients    map[*Client]bool
}

func newHub() *Hub {
	return &Hub{
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case msg := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- msg:
				default:
					close(client.send)
				}
			}
		}
	}
}

type Client struct {
	conn *websocket.Conn
	send chan []byte
	hub  *Hub
	log  *slog.Logger
}

func (c *Client) readPump() {
	const op = "http.ws.Client.readPump"
	log := c.log.With(slog.String("op", op)) // add user info

	defer func() {
		c.hub.unregister <- c
		_ = c.conn.Close()
	}()

	c.conn.SetReadLimit(512)
	heartbeat(c.conn)

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Error("failed to read message", logerr.Err(err))
			} else {
				log.Info("connection closed", slog.String("reason", err.Error()))
			}
			break
		}

		c.hub.broadcast <- msg
	}
}

func (c *Client) writePump() {
	const op = "http.ws.Client.writePump"
	log := c.log.With(slog.String("op", op)) // add user info

	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = c.conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := c.conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Error("failed to write message", logerr.Err(err))
				return
			}
		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			err := c.conn.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				c.log.Error("failed to write ping message", logerr.Err(err))
				return
			}
		}
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
		send: make(chan []byte, 256),
		hub:  h.hub,
		log:  h.log,
	}
	h.hub.register <- client

	go client.writePump()
	client.readPump()
}
