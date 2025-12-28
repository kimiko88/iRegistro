package ws

type Hub struct {
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
	Rooms      map[string]map[*Client]bool // RoomID -> Set of Clients
}

func NewHub() *Hub {
	return &Hub{
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Rooms:      make(map[string]map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
				h.removeClientFromRooms(client)
			}
		case message := <-h.Broadcast:
			for client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, client)
				}
			}
		}
	}
}

func (h *Hub) JoinRoom(room string, client *Client) {
	if h.Rooms[room] == nil {
		h.Rooms[room] = make(map[*Client]bool)
	}
	h.Rooms[room][client] = true
	client.Rooms[room] = true
}

func (h *Hub) removeClientFromRooms(client *Client) {
	for room := range client.Rooms {
		if clients, ok := h.Rooms[room]; ok {
			delete(clients, client)
			if len(clients) == 0 {
				delete(h.Rooms, room)
			}
		}
	}
}

func (h *Hub) BroadcastToRoom(room string, message []byte) {
	if clients, ok := h.Rooms[room]; ok {
		for client := range clients {
			select {
			case client.Send <- message:
			default:
				close(client.Send)
				delete(h.Clients, client)
			}
		}
	}
}
