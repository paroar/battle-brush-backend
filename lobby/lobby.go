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
	clients map[string]*Client
	rooms   map[string]IRoom
}

// NewLobby creates a Lobby
func NewLobby() *Lobby {
	return &Lobby{
		clients: make(map[string]*Client),
		rooms:   make(map[string]IRoom),
	}
}

// AddClient adds a Client to the Lobby
func (lobby *Lobby) AddClient(c *Client) {
	lobby.clients[c.id] = c
}

// DeleteClient deletes a Client from the Lobby
func (lobby *Lobby) DeleteClient(id string) {
	if c, ok := lobby.clients[id]; ok {
		delete(lobby.clients, c.id)
	}
}

// AddRoom adds a Room to the Lobby
func (lobby *Lobby) AddRoom(r IRoom) {
	lobby.rooms[r.GetID()] = r
}

// DeleteRoom deletes a Room from the Lobby
func (lobby *Lobby) DeleteRoom(id string) {
	if r, ok := lobby.rooms[id]; ok {
		delete(lobby.rooms, r.GetID())
	}
}

func (lobby *Lobby) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(rw, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := NewClient(lobby, conn)

	lobby.AddClient(client)

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
		room := lobby.GetRoom(roomid)
		if room != nil {
			pr := room.(IRoom)
			pr.JoinClient(client)
		}

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

// GetClient returns the Client if found or Error
func (lobby *Lobby) GetClient(id string) (*Client, error) {
	if c, ok := lobby.clients[id]; ok {
		return c, nil
	}

	return nil, errors.New("Client not found")
}

// GetRoom returns the Room if found or Error
func (lobby *Lobby) GetRoom(id string) interface{} {
	if r, ok := lobby.rooms[id]; ok {
		return r
	}

	return nil
}

// FirstAvailablePublicRoom returns an available public Room if there is one
func (lobby *Lobby) FirstAvailablePublicRoom() interface{} {
	for _, room := range lobby.rooms {
		//&& pr.game.state == StateWaiting
		if room.IsPublic() && room.IsAvailable() {
			return room
		}
	}
	return nil
}
