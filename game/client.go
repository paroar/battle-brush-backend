package game

import (
	"log"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/paroar/battle-brush-backend/message"
)

// Client struct
type Client struct {
	ID    string
	Conn  *websocket.Conn
	Lobby *Lobby
	Room  *Room
	Send  chan *message.Message
}

// NewClient creates a Client without a Room atached
func NewClient(lobby *Lobby) *Client {
	return &Client{
		ID:    uuid.NewString(),
		Lobby: lobby,
		Send:  make(chan *message.Message),
	}
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
		var msg message.Message
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			log.Println(err)
			return
		}
		if c.Room != nil {
			c.Room.Broadcast <- &msg
		}
	}
}

func (c *Client) disconnect() {
	if c.Room != nil {
		c.Room.LeaveClientChan <- c
	}
	c.Lobby.LeaveClientChan <- c
	c.Conn.Close()
}
