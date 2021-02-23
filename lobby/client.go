package lobby

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/paroar/battle-brush-backend/generators"
)

// Client struct
type Client struct {
	name  string
	id    string
	conn  *websocket.Conn
	lobby *Lobby
	room  *Room
	send  chan *Message
}

// NewClient creates a Client without a Room atached
func NewClient(lobby *Lobby, conn *websocket.Conn) *Client {
	return &Client{
		name:  generators.Name(),
		id:    uuid.NewString(),
		lobby: lobby,
		conn:  conn,
		send:  make(chan *Message),
	}
}

func (c *Client) writePump() {
	defer c.disconnect()

	for {
		msg := <-c.send

		err := c.conn.WriteJSON(msg)
		if err != nil {
			log.Println(err)
		}
	}
}

func (c *Client) readPump() {
	defer c.disconnect()

	for {
		var raw json.RawMessage
		msg := &Message{
			Content: &raw,
		}

		err := c.conn.ReadJSON(&msg)
		if err != nil {
			log.Println(err)
			return
		}

		content := TypeHandler[msg.Type]()
		if err := json.Unmarshal(raw, &content); err != nil {
			log.Println(err)
		}

		content.(IContent).Do(c)

	}
}

func (c *Client) disconnect() {
	if c.room != nil {
		c.room.leaveClientChan <- c
	}
	c.lobby.leaveClientChan <- c
	c.conn.Close()
}
