package ws

import (
	"github.com/dundduun/msg/core/pkg/logerr"
	"github.com/gorilla/websocket"
	"log/slog"
	"time"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = pongWait * 9 / 10
)

type Envelope struct {
	Type   string `json:"type"` // mandatory
	Text   string `json:"text"`
	Room   string `json:"room"`
	List   []Room `json:"list"`
	sender *Client
}

type Room struct {
	Name    string `json:"name"`
	Members int    `json:"members"`
}

type Client struct {
	room string
	send chan Envelope
	conn *websocket.Conn
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
		var env Envelope
		err := c.conn.ReadJSON(&env)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Error("failed to read message", logerr.Err(err))
				log.Info("connection closed")
			} else {
				log.Info("connection closed", slog.String("reason", err.Error()))
			}
			break
		}

		switch env.Type {
		case "join":
			c.hub.switchRooms(c, env.Room)
		case "leave":
			c.hub.switchRooms(c, "lobby")
		case "list":
			rooms := c.hub.countRooms()
			env = Envelope{
				Type:   "list_response",
				List:   rooms,
				sender: c,
			}
			c.hub.broadcast <- env
		case "broadcast": // rename response
			env.sender = c
			c.hub.broadcast <- env
		}
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
		case env, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := c.conn.WriteJSON(env)
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

func heartbeat(conn *websocket.Conn) {
	_ = conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		return conn.SetReadDeadline(time.Now().Add(pongWait))
	})
}
