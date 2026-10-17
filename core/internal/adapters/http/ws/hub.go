package ws

import (
	"github.com/google/uuid"
	"log/slog"
)

type Hub struct {
	log        *slog.Logger
	groups     map[string]map[*Client]bool
	users      map[uuid.UUID]*Client
	direct     chan Envelope
	broadcast  chan Envelope
	register   chan *Client
	unregister chan *Client
}

func newHub(log *slog.Logger) *Hub {
	hub := Hub{
		log:        log,
		groups:     make(map[string]map[*Client]bool),
		users:      make(map[uuid.UUID]*Client),
		direct:     make(chan Envelope),
		broadcast:  make(chan Envelope),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
	go hub.run()

	return &hub
}

func (h *Hub) run() {
	const op = "ws.Hub.run"
	log := h.log.With(slog.String("op", op))

	log.Debug("hub is running")

	for {
		log.Debug("for loop start")

		select {
		case c := <-h.register:
			log.Debug("reg start")
			if h.groups[c.group] == nil {
				h.groups[c.group] = make(map[*Client]bool)
			}
			h.groups[c.group][c] = true
			h.broadcastMessage(Envelope{
				Type:   "system.online",
				UserID: c.profile.ID,
				sender: c,
			})
			h.users[c.profile.ID] = c
		case c := <-h.unregister:
			log.Debug("start")
			if h.groups[c.group] != nil {
				delete(h.groups[c.group], c)
				if len(h.groups[c.group]) == 0 {
					delete(h.groups, c.group)
				}
				c.group = "lobby"
				close(c.send)
			}

			h.broadcastMessage(Envelope{
				Type:   "system.offline",
				UserID: c.profile.ID,
				sender: c,
			})
			delete(h.users, c.profile.ID)
			log.Debug("end")

		case env := <-h.direct:
			log.Debug("start")
			target, online := h.users[env.To]
			if !online {
				//h.queueOffline(env)
				continue
			}
			select {
			case target.send <- env:
			default:
				close(target.send)
			}
			log.Debug("end")

		case env := <-h.broadcast:
			h.broadcastMessage(env)
		}
	}
}

//func (h *Hub) queueOffline(Envelope) {} // implement

func (h *Hub) broadcastMessage(env Envelope) {
	const op = "ws.Hub.broadcastMessage"
	log := h.log.With(slog.String("op", op))

	log.Debug("broadcast case")
	for c := range h.groups[env.sender.group] {
		select {
		case c.send <- env:
		default:
			close(c.send)
		}
	}
}

func (h *Hub) switchGroups(c *Client, group string) {
	if h.groups[c.group] != nil {
		delete(h.groups[c.group], c)
		if len(h.groups[c.group]) == 0 {
			delete(h.groups, c.group)
		}
	}

	if h.groups[group] == nil {
		h.groups[group] = make(map[*Client]bool)
	}
	h.groups[group][c] = true
	c.group = group
}

func (h *Hub) listGroups() []Group {
	var groups []Group

	for name, group := range h.groups {
		groups = append(groups, Group{name, len(group)})
	}

	return groups
}

func (h *Hub) listOnline() []ProfileData {
	var profiles []ProfileData

	for _, client := range h.users {
		profiles = append(profiles, client.profile)
	}

	return profiles
}
