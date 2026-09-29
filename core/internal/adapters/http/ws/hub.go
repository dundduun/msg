package ws

type Hub struct {
	rooms      map[string]map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan Envelope
}

func newHub() *Hub {
	hub := Hub{
		rooms:      make(map[string]map[*Client]bool),
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
			if h.rooms[c.room] == nil {
				h.rooms[c.room] = make(map[*Client]bool)
			}
			h.rooms[c.room][c] = true

		case c := <-h.unregister:
			if h.rooms[c.room] != nil {
				delete(h.rooms[c.room], c)
				if len(h.rooms[c.room]) == 0 {
					delete(h.rooms, c.room)
				}
				c.room = "lobby"
				close(c.send)
			}

		case env := <-h.broadcast:
			for c := range h.rooms[env.sender.room] {
				select {
				case c.send <- env:
				default:
					close(c.send)
				}
			}
		}
	}
}

func (h *Hub) switchRooms(c *Client, room string) {
	if h.rooms[c.room] != nil {
		delete(h.rooms[c.room], c)
		if len(h.rooms[c.room]) == 0 {
			delete(h.rooms, c.room)
		}
	}

	if h.rooms[room] == nil {
		h.rooms[room] = make(map[*Client]bool)
	}
	h.rooms[room][c] = true
	c.room = room
}

func (h *Hub) countRooms() []Room {
	var rooms []Room

	for name, room := range h.rooms {
		rooms = append(rooms, Room{name, len(room)})
	}

	return rooms
}
