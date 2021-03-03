package websocket

import (
	"github.com/gorilla/websocket"
	"github.com/paroar/battle-brush-backend/message"
)

type Client struct {
	ID    string
	Conn  *websocket.Conn
	Lobby *Lobby
	Send  chan *message.Envelope
}

func NewClient(id string, conn *websocket.Conn, lobby *Lobby) *Client {
	return &Client{
		ID:    id,
		Conn:  conn,
		Lobby: lobby,
		Send:  make(chan *message.Envelope),
	}
}

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
			panic(err)
		}
	}
}

func (c *Client) readPump() {
	defer c.disconnect()

	for {

	}
}

func (c *Client) disconnect() {
	c.Conn.Close()
}
