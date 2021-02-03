package game

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Lobby is the main struct that manages rooms and clients
type Lobby struct {
	clients         map[*Client]bool
	JoinClientChan  chan *Client
	LeaveClientChan chan *Client
	rooms           map[*Room]bool
	JoinRoomChan    chan *Room
	LeaveRoomChan   chan *Room
	broadcast       chan []byte
}

// NewLobby creates a Lobby
func NewLobby() *Lobby {
	return &Lobby{
		clients:         make(map[*Client]bool),
		JoinClientChan:  make(chan *Client),
		LeaveClientChan: make(chan *Client),
		rooms:           make(map[*Room]bool),
		JoinRoomChan:    make(chan *Room),
		LeaveRoomChan:   make(chan *Room),
		broadcast:       make(chan []byte),
	}
}

func (lobby *Lobby) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(rw, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{
		ID:    uuid.NewString(),
		Lobby: lobby,
		Conn:  conn,
		Send:  make(chan []byte),
	}

	lobby.JoinClientChan <- client

	go client.writePump()
	go client.readPump()

	msg := &Message{
		Type:    "id",
		Content: client.ID,
	}
	_msg, err := json.Marshal(msg)
	client.Send <- _msg

}

// Run runs the lobby
func (lobby *Lobby) Run() {
	for {
		select {
		case client := <-lobby.JoinClientChan:
			lobby.joinClient(client)
		case client := <-lobby.LeaveClientChan:
			lobby.leaveClient(client)
		case room := <-lobby.JoinRoomChan:
			lobby.joinRoom(room)
		case room := <-lobby.LeaveRoomChan:
			lobby.leaveRoom(room)
		case msg := <-lobby.broadcast:
			lobby.broadcastTo(msg)
		}
	}
}

func (lobby *Lobby) joinClient(c *Client) {
	lobby.clients[c] = true
}

func (lobby *Lobby) leaveClient(c *Client) {
	if _, ok := lobby.clients[c]; ok {
		delete(lobby.clients, c)
	}
}

func (lobby *Lobby) joinRoom(r *Room) {
	lobby.rooms[r] = true
}

func (lobby *Lobby) leaveRoom(r *Room) {
	if _, ok := lobby.rooms[r]; ok {
		delete(lobby.rooms, r)
	}
}

func (lobby *Lobby) broadcastTo(msg []byte) {
	for client := range lobby.clients {
		client.Send <- msg
	}
}

// GetLobbyClient returns the Client if found or Error
func (lobby *Lobby) GetLobbyClient(id string) (*Client, error) {
	for client := range lobby.clients {
		if client.ID == id {
			return client, nil
		}
	}
	return nil, errors.New("Client not found")
}

// GetLobbyRoom returns the Room if found or Error
func (lobby *Lobby) GetLobbyRoom(id string) (*Room, error) {
	for room := range lobby.rooms {
		if room.ID == id {
			return room, nil
		}
	}
	return nil, errors.New("Room not found")
}

// GetLobbyRooms returns all Rooms in the lobby
func (lobby *Lobby) GetLobbyRooms() []Room {
	var rooms = []Room{}
	for room := range lobby.rooms {
		rooms = append(rooms, *room)
	}
	return rooms
}

// GetLobbyClients returns all Clients in the lobby
func (lobby *Lobby) GetLobbyClients() []Client {
	var clients = []Client{}
	for client := range lobby.clients {
		clients = append(clients, *client)
	}
	return clients
}
