package lobby

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// Room struct
type Room struct {
	lobby           *Lobby
	Clients         map[*Client]bool
	JoinClientChan  chan *Client
	LeaveClientChan chan *Client
	ID              string
	Broadcast       chan *Message
	Options         RoomOptions
	Game            *DrawGame
}

// RoomOptions struct
type RoomOptions struct {
	NumPlayers int
}

var defaultOptions = &RoomOptions{
	NumPlayers: 5,
}

// NewDefaultRoom creates a Room
func NewDefaultRoom(lobby *Lobby) *Room {
	clients := make(map[*Client]bool)
	return &Room{
		lobby:           lobby,
		Clients:         clients,
		JoinClientChan:  make(chan *Client),
		LeaveClientChan: make(chan *Client),
		ID:              uuid.NewString(),
		Broadcast:       make(chan *Message),
		Options:         *defaultOptions,
		Game:            NewDrawGame(clients),
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
		Broadcast:       make(chan *Message),
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
	_msg := &Message{
		Type: TypeJoinLeave,
		Content: JoinLeave{
			UserName: c.Name,
			ID:       c.ID,
			Msg:      fmt.Sprintf("%s has joined", c.Name),
		},
	}
	room.broadcastTo(_msg)

	_msg = &Message{
		Type: TypePlayers,
		Content: Players{
			UserNames: room.getUserNames(),
		},
	}
	room.broadcastTo(_msg)
}

func (room *Room) leaveClient(c *Client) {
	if _, ok := room.Clients[c]; ok {
		delete(room.Clients, c)
		_msg := &Message{
			Type: TypeJoinLeave,
			Content: JoinLeave{
				UserName: c.Name,
				ID:       c.ID,
				Msg:      fmt.Sprintf("%s has left", c.Name),
			},
		}
		room.broadcastTo(_msg)
		_msg = &Message{
			Type: TypePlayers,
			Content: Players{
				UserNames: room.getUserNames(),
			},
		}
		room.broadcastTo(_msg)
	}
}

func (room *Room) broadcastTo(msg *Message) {
	for client := range room.Clients {
		client.Send <- msg
	}
}

func (room *Room) getClient(userid string) (*Client, error) {
	for client := range room.Clients {
		if client.ID == userid {
			return client, nil
		}
	}
	return nil, errors.New("Client not found")
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
