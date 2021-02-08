package game

import (
	"github.com/google/uuid"
	"github.com/paroar/battle-brush-backend/message"
)

// Room struct
type Room struct {
	lobby           *Lobby
	Clients         map[*Client]bool
	JoinClientChan  chan *Client
	LeaveClientChan chan *Client
	theme           string
	ID              string
	Broadcast       chan *message.Message
	RoomOptions     RoomOptions
}

// RoomOptions struct
type RoomOptions struct {
	NumPlayers int
	Time       int
	Rounds     int
	Private    bool
}

// NewRoom creates a Room
func NewRoom(lobby *Lobby, roomOptions RoomOptions) *Room {
	return &Room{
		lobby:           lobby,
		Clients:         make(map[*Client]bool),
		JoinClientChan:  make(chan *Client),
		LeaveClientChan: make(chan *Client),
		theme:           "beach",
		ID:              uuid.NewString(),
		Broadcast:       make(chan *message.Message),
		RoomOptions:     roomOptions,
	}
}

// Run runs the Room
func (r *Room) Run() {
	for {
		select {
		case client := <-r.JoinClientChan:
			r.joinClient(client)
		case client := <-r.LeaveClientChan:
			r.leaveClient(client)
		case msg := <-r.Broadcast:
			r.broadcastTo(msg)
		}
	}
}

func (r *Room) joinClient(c *Client) {
	r.Clients[c] = true
	_msg := &message.Message{
		Type: message.TypeJoin,
		Content: message.Join{
			ID: c.ID,
		},
	}
	r.broadcastTo(_msg)
}

func (r *Room) leaveClient(c *Client) {
	if _, ok := r.Clients[c]; ok {
		delete(r.Clients, c)
		_msg := &message.Message{
			Type: message.TypeLeave,
			Content: message.Leave{
				ID: c.ID,
			},
		}
		r.broadcastTo(_msg)
	}
}

func (r *Room) broadcastTo(msg *message.Message) {
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
