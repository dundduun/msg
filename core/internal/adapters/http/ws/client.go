package ws

import (
	"github.com/dundduun/msg/core/pkg/logerr"
	"github.com/google/uuid"
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
	Type       string        `json:"type"`
	Text       string        `json:"text,omitempty"`
	UserID     uuid.UUID     `json:"user-id,omitzero"`
	Group      string        `json:"group,omitempty"`
	GroupList  []Group       `json:"group-list,omitempty"`
	OnlineList []ProfileData `json:"online-list,omitempty"`
	To         uuid.UUID     `json:"to,omitzero"`
	sender     *Client
}

type Group struct {
	Name    string `json:"name"`
	Members int    `json:"members"`
}

type Client struct {
	clientID uuid.UUID
	profile  ProfileData
	group    string
	send     chan Envelope
	conn     *websocket.Conn
	hub      *Hub
	log      *slog.Logger
}

type ProfileData struct {
	ID       uuid.UUID `json:"id,omitempty"`
	Name     string    `json:"name,omitempty"`
	Username string    `json:"username,omitempty"`
}

func (c *Client) readPump() {
	const op = "http.ws.Client.readPump"
	log := c.log.With(slog.String("op", op), slog.String("id", c.profile.ID.String()))

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
		log.Debug("new message", slog.Any("env", env))

		switch env.Type {
		case "join":
			c.hub.switchGroups(c, env.Group)
		case "leave":
			c.hub.switchGroups(c, "lobby")
		case "list_groups":
			groups := c.hub.listGroups()
			env = Envelope{
				Type:      "list_groups_response",
				GroupList: groups,
				sender:    c,
			}
			c.hub.broadcast <- env
		case "list_online":
			online := c.hub.listOnline()
			env = Envelope{
				Type:       "list_online_response",
				OnlineList: online,
				sender:     c,
			}
			c.hub.broadcast <- env
		case "direct":
			env.sender = c
			c.hub.direct <- env
		case "broadcast": // rename response
			log.Debug("broadcasting")
			env.sender = c
			env.Type = "broadcast_response"
			c.hub.broadcast <- env
			log.Debug("broadcast env sent to chan")
		}
	}
}

func (c *Client) writePump() {
	const op = "http.ws.Client.writePump"
	log := c.log.With(slog.String("op", op), slog.String("id", c.profile.ID.String()))

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
