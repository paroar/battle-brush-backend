package game

import (
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/paroar/battle-brush-backend/message"
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
	clients              map[*Client]bool
	JoinClientChan       chan *Client
	LeaveClientChan      chan *Client
	rooms                *Rooms
	JoinPublicRoomChan   chan *Room
	LeavePublicRoomChan  chan *Room
	JoinPrivateRoomChan  chan *Room
	LeavePrivateRoomChan chan *Room
	broadcast            chan *message.Message
}

// NewLobby creates a Lobby
func NewLobby() *Lobby {
	return &Lobby{
		clients:              make(map[*Client]bool),
		JoinClientChan:       make(chan *Client),
		LeaveClientChan:      make(chan *Client),
		rooms:                NewRooms(),
		JoinPublicRoomChan:   make(chan *Room),
		LeavePublicRoomChan:  make(chan *Room),
		JoinPrivateRoomChan:  make(chan *Room),
		LeavePrivateRoomChan: make(chan *Room),
		broadcast:            make(chan *message.Message),
	}
}

func (lobby *Lobby) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(rw, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := NewClient(lobby, conn)

	lobby.JoinClientChan <- client

	go client.readPump()
	go client.writePump()

	_msg := &message.Message{
		Type: message.TypeLogin,
		Content: message.Login{
			UserName: client.Name,
			ID:       client.ID,
		},
	}

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
		case room := <-lobby.JoinPublicRoomChan:
			lobby.joinPublicRoom(room)
		case room := <-lobby.LeavePublicRoomChan:
			lobby.leavePublicRoom(room)
		case room := <-lobby.JoinPrivateRoomChan:
			lobby.joinPrivateRoom(room)
		case room := <-lobby.LeavePrivateRoomChan:
			lobby.leavePrivateRoom(room)
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

func (lobby *Lobby) joinPublicRoom(r *Room) {
	lobby.rooms.publicRooms[r] = true
}

func (lobby *Lobby) leavePublicRoom(r *Room) {
	if _, ok := lobby.rooms.publicRooms[r]; ok {
		delete(lobby.rooms.publicRooms, r)
	}
}

func (lobby *Lobby) joinPrivateRoom(r *Room) {
	lobby.rooms.privateRooms[r] = true
}

func (lobby *Lobby) leavePrivateRoom(r *Room) {
	if _, ok := lobby.rooms.privateRooms[r]; ok {
		delete(lobby.rooms.privateRooms, r)
	}
}

func (lobby *Lobby) broadcastTo(msg *message.Message) {
	for client := range lobby.clients {
		client.Send <- msg
	}
}

// GetClient returns the Client if found or Error
func (lobby *Lobby) GetClient(id string) (*Client, error) {
	for client := range lobby.clients {
		if client.ID == id {
			return client, nil
		}
	}
	return nil, errors.New("Client not found")
}

// GetPublicRoom returns the Room if found or Error
func (lobby *Lobby) GetPublicRoom(id string) (*Room, error) {
	room, err := lobby.rooms.GetPublicRoom(id)
	if err != nil {
		return nil, err
	}
	return room, nil
}

// GetPrivateRoom returns the Room if found or Error
func (lobby *Lobby) GetPrivateRoom(id string) (*Room, error) {
	room, err := lobby.rooms.GetPrivateRoom(id)
	if err != nil {
		return nil, err
	}
	return room, nil
}

// GetPublicRooms returns all Rooms in the lobby
func (lobby *Lobby) GetPublicRooms() []Room {
	var rooms = []Room{}
	for room := range lobby.rooms.publicRooms {
		rooms = append(rooms, *room)
	}
	return rooms
}

// GetClients returns all Clients in the lobby
func (lobby *Lobby) GetClients() []Client {
	var clients = []Client{}
	for client := range lobby.clients {
		clients = append(clients, *client)
	}
	return clients
}

// CreatePrivateRoom creates the Room
func (lobby *Lobby) CreatePrivateRoom(roomOptions *RoomOptions, client *Client) *Room {
	room := NewPrivateRoom(lobby, roomOptions)
	go room.Run()
	lobby.JoinPrivateRoomChan <- room
	lobby.LeaveClientChan <- client
	room.JoinClientChan <- client
	client.Room = room

	return room
}

// JoinPrivateRoom returns an error if the Room is full or joins the Room
func (lobby *Lobby) JoinPrivateRoom(room *Room, client *Client) error {

	if room.isFull() {
		return errors.New("Room is full")
	}

	lobby.LeaveClientChan <- client
	room.JoinClientChan <- client
	client.Room = room

	return nil
}

// CreateOrJoinPublicRoom creates a Room if there are no available
// public Rooms or joins a Room if there is one available
func (lobby *Lobby) CreateOrJoinPublicRoom(client *Client) *Room {
	room := lobby.AvailablePublicRooms()
	if room == nil {
		room = NewDefaultRoom(lobby)
		go room.Run()
		lobby.JoinPublicRoomChan <- room
	}

	lobby.LeaveClientChan <- client
	room.JoinClientChan <- client
	client.Room = room

	return room
}

// AvailablePublicRooms returns an available public Room if there is one
func (lobby *Lobby) AvailablePublicRooms() *Room {
	for room := range lobby.rooms.publicRooms {
		if len(room.Clients) <= room.Options.NumPlayers {
			return room
		}
	}
	return nil
}
