package lobby

import (
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
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
	clients              map[*Client]bool
	joinClientChan       chan *Client
	leaveClientChan      chan *Client
	rooms                *Rooms
	joinPublicRoomChan   chan *Room
	leavePublicRoomChan  chan *Room
	joinPrivateRoomChan  chan *Room
	leavePrivateRoomChan chan *Room
	broadcast            chan *Message
}

// NewLobby creates a Lobby
func NewLobby() *Lobby {
	return &Lobby{
		clients:              make(map[*Client]bool),
		joinClientChan:       make(chan *Client),
		leaveClientChan:      make(chan *Client),
		rooms:                NewRooms(),
		joinPublicRoomChan:   make(chan *Room),
		leavePublicRoomChan:  make(chan *Room),
		joinPrivateRoomChan:  make(chan *Room),
		leavePrivateRoomChan: make(chan *Room),
		broadcast:            make(chan *Message),
	}
}

func (lobby *Lobby) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(rw, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := NewClient(lobby, conn)

	lobby.joinClientChan <- client

	go client.readPump()
	go client.writePump()

	msg := &Message{
		Type: TypeLogin,
		Content: Login{
			UserName: client.name,
			ID:       client.id,
		},
	}

	client.send <- msg

	vars := mux.Vars(r)
	roomid := vars["room"]

	if roomid != "" {
		client, err := lobby.GetClient(client.id)
		room, err := lobby.GetPrivateRoom(roomid)
		err = lobby.JoinClientToPrivateRoom(room, client)

		msg := &Message{
			Type: TypeConnection,
		}

		if err != nil {
			msg.Content = Connection{
				RoomID: roomid,
				Status: err.Error(),
			}
		} else {
			msg.Content = Connection{
				RoomID:   roomid,
				Status:   "ok",
				RoomType: RoomTypePrivate,
			}
		}

		client.send <- msg
	}

}

// Run runs the lobby
func (lobby *Lobby) Run() {
	for {
		select {
		case client := <-lobby.joinClientChan:
			lobby.joinClient(client)
		case client := <-lobby.leaveClientChan:
			lobby.leaveClient(client)
		case room := <-lobby.joinPublicRoomChan:
			lobby.joinPublicRoom(room)
		case room := <-lobby.leavePublicRoomChan:
			lobby.leavePublicRoom(room)
		case room := <-lobby.joinPrivateRoomChan:
			lobby.joinPrivateRoom(room)
		case room := <-lobby.leavePrivateRoomChan:
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

func (lobby *Lobby) broadcastTo(msg *Message) {
	for client := range lobby.clients {
		client.send <- msg
	}
}

// GetClient returns the Client if found or Error
func (lobby *Lobby) GetClient(id string) (*Client, error) {
	for client := range lobby.clients {
		if client.id == id {
			return client, nil
		}
	}
	return nil, errors.New("Client not found")
}

// GetPublicRoom returns the Room if found or Error
func (lobby *Lobby) GetPublicRoom(id string) (*Room, error) {
	return lobby.rooms.GetPublicRoom(id)
}

// GetPrivateRoom returns the Room if found or Error
func (lobby *Lobby) GetPrivateRoom(id string) (*Room, error) {
	return lobby.rooms.GetPrivateRoom(id)
}

// CreatePrivateRoom creates the Room
func (lobby *Lobby) CreatePrivateRoom(client *Client) *Room {
	room := NewPrivateRoom(lobby)
	go room.run()
	go room.game.run()
	client.room = room
	lobby.joinPrivateRoomChan <- room
	lobby.leaveClientChan <- client
	room.joinClientChan <- client

	return room
}

// JoinClientToPrivateRoom returns an error if the Room is full or joins the Room
func (lobby *Lobby) JoinClientToPrivateRoom(room *Room, client *Client) error {

	if room.isAvailable() {
		return errors.New("Room is full")
	}

	client.room = room
	lobby.leaveClientChan <- client
	room.joinClientChan <- client

	return nil
}

// CreateOrJoinPublicRoom creates a Room if there are no available
// public Rooms or joins a Room if there is one available
func (lobby *Lobby) CreateOrJoinPublicRoom(client *Client) *Room {
	room := lobby.firstAvailablePublicRoom()
	if room == nil {
		room = NewDefaultRoom(lobby)
		go room.run()
		go room.game.run()
		lobby.joinPublicRoomChan <- room
	}

	client.room = room
	lobby.leaveClientChan <- client
	room.joinClientChan <- client

	return room
}

// firstAvailablePublicRooms returns an available public Room if there is one
func (lobby *Lobby) firstAvailablePublicRoom() *Room {
	for room := range lobby.rooms.publicRooms {
		if room.game.state == StateWaiting && len(room.clients) <= room.options.NumPlayers {
			return room
		}
	}
	return nil
}

// AvailableRooms returns all available Rooms
func (lobby *Lobby) availableRooms(rooms []*Room) []*Room {
	availableRooms := []*Room{}
	for _, room := range rooms {
		if room.game.state == StateWaiting && len(room.clients) <= room.options.NumPlayers {
			availableRooms = append(availableRooms, room)
		}
	}
	return availableRooms
}
