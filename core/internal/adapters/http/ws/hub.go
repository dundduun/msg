package ws

type Hub struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan Envelope
}

func newHub() *Hub {
	hub := Hub{
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan Envelope),
		clients:    make(map[*Client]bool),
	}
	go hub.run()

	return &hub
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
