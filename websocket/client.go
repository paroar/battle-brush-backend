package websocket

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/paroar/battle-brush-backend/db"
	"github.com/paroar/battle-brush-backend/message"
	"github.com/paroar/battle-brush-backend/message/content"
)

// Client struct
type Client struct {
	ID    string
	Conn  *websocket.Conn
	Lobby *Lobby
	Send  chan *message.Envelope
}

// NewClient contructor
func NewClient(id string, conn *websocket.Conn, lobby *Lobby) *Client {
	return &Client{
		ID:    id,
		Conn:  conn,
		Lobby: lobby,
		Send:  make(chan *message.Envelope),
	}
}

// Run gets the websocket connection running in parallel
func (c *Client) Run() {
	go c.readPump()
	go c.writePump()
}

func (c *Client) writePump() {
	defer c.disconnect()

	for {
		msg := <-c.Send

		err := c.Conn.WriteJSON(msg)
		if err != nil {
			log.Println(err)
		}
	}
}

func (c *Client) readPump() {
	defer c.disconnect()

	for {
		var raw json.RawMessage
		msg := &message.Envelope{
			Content: &raw,
		}

		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func (c *Client) disconnect() {
	defer c.Conn.Close()

	player, err := db.DeletePlayer(c.ID)
	if err != nil {
		log.Println(err)
		return
	}
	delete(c.Lobby.Clients, c.ID)

	room, err := db.ReadRoom(player.RoomID)
	if err != nil {
		log.Println(err)
		return
	}

	playersNames := db.ReadPlayersNames(room.PlayersID)

	msg := &message.Envelope{
		Type: content.TypePlayers,
		Content: content.Players{
			UserNames: playersNames,
		},
	}
	c.Lobby.Broadcast(room.PlayersID, msg)

	msg = &message.Envelope{
		Type: content.TypeJoinLeave,
		Content: content.JoinLeave{
			UserName: player.Name,
			ID:       player.ID,
			Msg:      fmt.Sprintf("%s has left", player.Name),
		},
	}
	c.Lobby.Broadcast(room.PlayersID, msg)

}
