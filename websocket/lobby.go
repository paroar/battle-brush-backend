package websocket

import (
	"fmt"
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

	msg := &message.Envelope{
		Type: content.TypeLogin,
		Content: content.Login{
			UserName: player.Name,
			ID:       player.ID,
		},
	}

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

	msg := &message.Envelope{
		Type: message.TypeConnection,
	}

	if err != nil {
		msg.Content = content.Connection{
			RoomID: roomid,
			Status: "room not found",
		}
	} else {
		msg.Content = content.Connection{
			RoomID:   roomid,
			Status:   "ok",
			RoomType: model.RoomTypePrivate,
		}
	}

	client.Send <- msg

	updatedPlayers := append(room.PlayersID, player.ID)
	room.UpdateRoom(updatedPlayers, room.State)
	db.UpdateRoom(room)

	player.RoomID = room.ID
	db.UpdatePlayer(player)

	playersNames := db.ReadPlayersNames(room.PlayersID)
	msg = &message.Envelope{
		Type: content.TypePlayers,
		Content: content.Players{
			UserNames: playersNames,
		},
	}
	l.Broadcast(room.PlayersID, msg)

	msg = &message.Envelope{
		Type: content.TypeJoinLeave,
		Content: content.JoinLeave{
			UserName: player.Name,
			ID:       player.ID,
			Msg:      fmt.Sprintf("%s has joined", player.Name),
		},
	}
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
