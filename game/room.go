package game

import "github.com/google/uuid"

// Room struct
type Room struct {
	lobby           *Lobby
	Clients         map[*Client]bool
	JoinClientChan  chan *Client
	LeaveClientChan chan *Client
	theme           string
	ID              string
	broadcast       chan []byte
}

// NewRoom creates a Room
func NewRoom(lobby *Lobby) *Room {
	return &Room{
		lobby:           lobby,
		Clients:         make(map[*Client]bool),
		JoinClientChan:  make(chan *Client),
		LeaveClientChan: make(chan *Client),
		theme:           "beach",
		ID:              uuid.NewString(),
		broadcast:       make(chan []byte),
	}
}

// Run runs the Room
func (r *Room) Run() {
	for {
		select {
		case client := <-r.JoinClientChan:
			r.registerClient(client)
		case client := <-r.LeaveClientChan:
			r.unregisterClient(client)
		case msg := <-r.broadcast:
			r.broadcastToClients(msg)
		}
	}
}

func (r *Room) registerClient(c *Client) {
	r.Clients[c] = true
}

func (r *Room) unregisterClient(c *Client) {
	delete(r.Clients, c)
}

func (r *Room) broadcastToClients(msg []byte) {
	for client := range r.Clients {
		client.Send <- msg
	}
}

// GetRoomClients returns all Clients in the Room
func GetRoomClients(r *Room) []Client {
	var clients = []Client{}
	for client := range r.Clients {
		clients = append(clients, *client)
	}
	return clients
}

// func CreateRoom(lobby *Lobby, client *Client) {

// 	room := NewRoom(lobby)
// 	room.Run()

// 	lobby.JoinRoomChan <- room
// 	room.JoinClientChan <- client

// }
