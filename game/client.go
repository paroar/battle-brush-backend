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

// NewClientWithID creates a Client without a Room atached and an ID
func NewClientWithID(lobby *Lobby, ID string) *Client {
	return &Client{
		ID:    ID,
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

		json.Marshal(_msg)
		err = c.Conn.WriteJSON(_msg)
		if err != nil {
			log.Println(err)
		}
	}
}

func (c *Client) readPump() {
	defer c.disconnect()

	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		c.Lobby.broadcast <- msg
	}
}

func (c *Client) disconnect() {
	c.Lobby.LeaveClientChan <- c
	c.Room.LeaveClientChan <- c
	c.Conn.Close()
}
