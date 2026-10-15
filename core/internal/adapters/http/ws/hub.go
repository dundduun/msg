package ws

import "github.com/google/uuid"

type Hub struct {
	groups     map[string]map[*Client]bool
	users      map[uuid.UUID]*Client
	direct     chan Envelope
	broadcast  chan Envelope
	register   chan *Client
	unregister chan *Client
}

func newHub() *Hub {
	hub := Hub{
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
	for {
		select {
		case c := <-h.register:
			if h.groups[c.group] == nil {
				h.groups[c.group] = make(map[*Client]bool)
			}
			h.groups[c.group][c] = true
			h.broadcast <- Envelope{
				Type:   "system.online",
				UserID: c.id,
			}
			h.users[c.id] = c

		case c := <-h.unregister:
			if h.groups[c.group] != nil {
				delete(h.groups[c.group], c)
				if len(h.groups[c.group]) == 0 {
					delete(h.groups, c.group)
				}
				c.group = "lobby"
				close(c.send)
			}
			h.broadcast <- Envelope{
				Type:   "system.offline",
				UserID: c.id,
			}
			delete(h.users, c.id)

		case env := <-h.direct:
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

		case env := <-h.broadcast:
			for c := range h.groups[env.sender.group] {
				select {
				case c.send <- env:
				default:
					close(c.send)
				}
			}
		}
	}
}

//func (h *Hub) queueOffline(Envelope) {} // implement

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

func (h *Hub) countGroups() []Group {
	var groups []Group

	for name, group := range h.groups {
		groups = append(groups, Group{name, len(group)})
	}

	return groups
}
