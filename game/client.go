package game

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Client struct
type Client struct {
	ID    string
	Conn  *websocket.Conn
	Lobby *Lobby
	Room  *Room
	Send  chan []byte
}

// NewClient creates a Client without a Room atached
func NewClient(lobby *Lobby) *Client {
	return &Client{
		ID:    uuid.NewString(),
		Lobby: lobby,
		Send:  make(chan []byte),
	}
}

func (c *Client) writePump() {
	defer c.disconnect()

	for {
		msg := <-c.Send

		var _msg Message
		err := json.Unmarshal(msg, &_msg)
		if err != nil {
			log.Println(err)
		}

		err = c.Conn.WriteJSON(_msg)
		if err != nil {
			log.Println(err)
		}
	}
}

func (c *Client) readPump() {
	defer c.disconnect()

	for {
		var _msg Message
		err := c.Conn.ReadJSON(&_msg)
		if err != nil {
			log.Println(err)
			return
		}
		msg, err := json.Marshal(_msg)
		if err != nil {
			log.Println(err)
			return
		}
		if c.Room != nil {
			c.Room.Broadcast <- msg
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
