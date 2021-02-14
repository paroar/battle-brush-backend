package lobby

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
	ID              string
	Broadcast       chan *message.Message
	Options         RoomOptions
}

// RoomOptions struct
type RoomOptions struct {
	NumPlayers int
	Time       int
	Rounds     int
}

var defaultOptions = &RoomOptions{
	NumPlayers: 5,
	Time:       90,
	Rounds:     3,
}

// NewDefaultRoom creates a Room
func NewDefaultRoom(lobby *Lobby) *Room {
	return &Room{
		lobby:           lobby,
		Clients:         make(map[*Client]bool),
		JoinClientChan:  make(chan *Client),
		LeaveClientChan: make(chan *Client),
		ID:              uuid.NewString(),
		Broadcast:       make(chan *message.Message),
		Options:         *defaultOptions,
	}
}

// NewPrivateRoom creates a Room
func NewPrivateRoom(lobby *Lobby, roomOptions *RoomOptions) *Room {
	return &Room{
		lobby:           lobby,
		Clients:         make(map[*Client]bool),
		JoinClientChan:  make(chan *Client),
		LeaveClientChan: make(chan *Client),
		ID:              uuid.NewString(),
		Broadcast:       make(chan *message.Message),
		Options:         *roomOptions,
	}
}

// Run runs the Room
func (room *Room) Run() {
	for {
		select {
		case client := <-room.JoinClientChan:
			room.joinClient(client)
		case client := <-room.LeaveClientChan:
			room.leaveClient(client)
		case msg := <-room.Broadcast:
			room.broadcastTo(msg)
		}
	}
}

func (room *Room) joinClient(c *Client) {
	room.Clients[c] = true
	_msg := &message.Message{
		Type: message.TypeJoin,
		Content: message.Join{
			UserName: c.Name,
			ID:       c.ID,
		},
	}
	room.broadcastTo(_msg)

	_msg = &message.Message{
		Type: message.TypePlayers,
		Content: message.Players{
			UserNames: room.getUserNames(),
		},
	}
	room.broadcastTo(_msg)
}

func (room *Room) leaveClient(c *Client) {
	if _, ok := room.Clients[c]; ok {
		delete(room.Clients, c)
		_msg := &message.Message{
			Type: message.TypeLeave,
			Content: message.Leave{
				UserName: c.Name,
				ID:       c.ID,
			},
		}
		room.broadcastTo(_msg)
		_msg = &message.Message{
			Type: message.TypePlayers,
			Content: message.Players{
				UserNames: room.getUserNames(),
			},
		}
		room.broadcastTo(_msg)
	}
}

func (room *Room) broadcastTo(msg *message.Message) {
	for client := range room.Clients {
		client.Send <- msg
	}
}

func (room *Room) getUserNames() []string {
	var userNames []string
	for client := range room.Clients {
		userNames = append(userNames, client.Name)
	}
	return userNames
}

// GetRoomClients returns all Clients in the Room
func GetRoomClients(room *Room) []Client {
	var clients = []Client{}
	for client := range room.Clients {
		clients = append(clients, *client)
	}
	return clients
}

func (room *Room) isFull() bool {
	return len(room.Clients) >= room.Options.NumPlayers
}
