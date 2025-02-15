package websocket

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/paroar/battle-brush-backend/db"
	"github.com/paroar/battle-brush-backend/message"
	"github.com/paroar/battle-brush-backend/message/content"
	"github.com/paroar/battle-brush-backend/model"
)

// Lobby struct
type Lobby struct {
	Clients map[string]*Client
}

// NewLobby constructor
func NewLobby() *Lobby {
	return &Lobby{
		Clients: make(map[string]*Client),
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (l *Lobby) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(rw, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	player := model.NewPlayer()
	db.CreatePlayer(player)
	client := NewClient(player.ID, conn, l)
	l.Clients[player.ID] = client
	client.Run()

	msg := content.NewLogin(player)
	client.Send <- msg

	vars := mux.Vars(r)
	roomid := vars["room"]

	if roomid != "" {
		l.joinByURL(roomid, client, player)
	}
}

func (l *Lobby) joinByURL(roomid string, client *Client, player *model.Player) {

	room, err := db.ReadRoom(roomid)
	if err != nil {
		log.Println(err)
	}

	if room.State != "Waiting" {
		return
	}

	var msg *message.Envelope
	if err != nil {
		msg = content.NewConnection(roomid, "room not found", "")
	} else {
		msg = content.NewConnection(roomid, "ok", model.RoomTypePrivate)
	}

	client.Send <- msg

	updatedPlayers := append(room.PlayersID, player.ID)
	room.UpdateRoom(updatedPlayers, room.State)
	db.UpdateRoom(room)

	player.RoomID = room.ID
	db.UpdatePlayer(player)

	playersNames := db.ReadPlayers(room.PlayersID)
	msg = content.NewPlayers(playersNames)
	l.Broadcast(room.PlayersID, msg)

	msg = content.NewJoinLeave(player, "has joined")
	l.Broadcast(room.PlayersID, msg)
}

// Broadcast a message to players
func (l *Lobby) Broadcast(playersid []string, msg *message.Envelope) {
	for _, player := range playersid {
		if client := l.Clients[player]; client != nil {
			client.Send <- msg
		}
	}
}
