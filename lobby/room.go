package lobby

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// Room struct
type Room struct {
	lobby           *Lobby
	clients         map[*Client]bool
	joinClientChan  chan *Client
	leaveClientChan chan *Client
	ID              string
	broadcast       chan *Message
	options         RoomOptions
	game            *DrawGame
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
		clients:         clients,
		joinClientChan:  make(chan *Client),
		leaveClientChan: make(chan *Client),
		ID:              uuid.NewString(),
		broadcast:       make(chan *Message),
		options:         *defaultOptions,
		game:            NewDrawGame(clients),
	}
}

// NewPrivateRoom creates a Room
func NewPrivateRoom(lobby *Lobby, roomOptions *RoomOptions) *Room {
	clients := make(map[*Client]bool)
	return &Room{
		lobby:           lobby,
		clients:         clients,
		joinClientChan:  make(chan *Client),
		leaveClientChan: make(chan *Client),
		ID:              uuid.NewString(),
		broadcast:       make(chan *Message),
		options:         *defaultOptions,
		game:            NewDrawGame(clients),
	}
}

// run runs the Room
func (room *Room) run() {
	for {
		select {
		case client := <-room.joinClientChan:
			room.joinClient(client)
		case client := <-room.leaveClientChan:
			room.leaveClient(client)
		case msg := <-room.broadcast:
			room.broadcastTo(msg)
		}
	}
}

func (room *Room) joinClient(c *Client) {
	room.clients[c] = true
	msg := &Message{
		Type: TypeJoinLeave,
		Content: JoinLeave{
			UserName: c.name,
			ID:       c.id,
			Msg:      fmt.Sprintf("%s has joined", c.name),
		},
	}
	room.broadcastTo(msg)

	msg = &Message{
		Type: TypePlayers,
		Content: Players{
			UserNames: room.getUserNames(),
		},
	}
	room.broadcastTo(msg)
}

func (room *Room) leaveClient(c *Client) {
	if _, ok := room.clients[c]; ok {
		delete(room.clients, c)
		msg := &Message{
			Type: TypeJoinLeave,
			Content: JoinLeave{
				UserName: c.name,
				ID:       c.id,
				Msg:      fmt.Sprintf("%s has left", c.name),
			},
		}
		room.broadcastTo(msg)
		msg = &Message{
			Type: TypePlayers,
			Content: Players{
				UserNames: room.getUserNames(),
			},
		}
		room.broadcastTo(msg)
	}
}

func (room *Room) broadcastTo(msg *Message) {
	for client := range room.clients {
		client.send <- msg
	}
}

func (room *Room) getClient(userid string) (*Client, error) {
	for client := range room.clients {
		if client.id == userid {
			return client, nil
		}
	}
	return nil, errors.New("Client not found")
}

func (room *Room) getUserNames() []string {
	var userNames []string
	for client := range room.clients {
		userNames = append(userNames, client.name)
	}
	return userNames
}

func (room *Room) isAvailable() bool {
	return room.game.state == StateWaiting && len(room.clients) >= room.options.NumPlayers
}
