package game

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/paroar/battle-brush-backend/generators"
	"github.com/paroar/battle-brush-backend/message"
)

// Client struct
type Client struct {
	Name  string
	ID    string
	Conn  *websocket.Conn
	Lobby *Lobby
	Room  *Room
	Send  chan *message.Message
}

// NewClient creates a Client without a Room atached
func NewClient(lobby *Lobby, conn *websocket.Conn) *Client {
	return &Client{
		Name:  generators.Name(),
		ID:    uuid.NewString(),
		Lobby: lobby,
		Conn:  conn,
		Send:  make(chan *message.Message),
	}
}

func (c *Client) writePump() {
	defer c.disconnect()

	for {
		_msg := <-c.Send

		err := c.Conn.WriteJSON(_msg)
		if err != nil {
			log.Println(err)
		}
	}
}

func (c *Client) readPump() {
	defer c.disconnect()

	for {
		var content json.RawMessage
		msg := message.Message{
			Content: &content,
		}
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			log.Println(err)
			return
		}
		switch msg.Type {
		case "Login":
			var l message.Login
			if err := json.Unmarshal(content, &l); err != nil {
				log.Println(err)
			}
		case "Chat":
			var c message.Chat
			if err := json.Unmarshal(content, &c); err != nil {
				log.Println(err)
			}
		case "Join":
			var j message.Join
			if err := json.Unmarshal(content, &j); err != nil {
				log.Println(err)
			}
		case "Leave":
			var l message.Leave
			if err := json.Unmarshal(content, &l); err != nil {
				log.Println(err)
			}
		default:
			log.Printf("unknown message type: %s", msg.Type)
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
