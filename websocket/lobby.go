package websocket

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/paroar/battle-brush-backend/db"
	"github.com/paroar/battle-brush-backend/lobby"
	"github.com/paroar/battle-brush-backend/message"
	"github.com/paroar/battle-brush-backend/model"
)

type Lobby struct {
	Clients map[string]*Client
}

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
		Type: lobby.TypeLogin,
		Content: lobby.Login{
			UserName: player.Name,
			ID:       player.ID,
		},
	}

	client.Send <- msg

	vars := mux.Vars(r)
	roomid := vars["room"]

	if roomid != "" {
		l.joinByUrl(roomid, client)
	}
}

func (l *Lobby) joinByUrl(roomid string, client *Client) {
	players := db.ReadRoomPlayers(roomid)

	msg := &message.Envelope{
		Type: message.TypeConnection,
	}

	if players != nil {
		msg.Content = message.Connection{
			RoomID:   roomid,
			Status:   "ok",
			RoomType: message.RoomTypePrivate,
		}
	} else {
		msg.Content = message.Connection{
			RoomID: roomid,
			Status: "room not found",
		}
	}

	client.Send <- msg
}
