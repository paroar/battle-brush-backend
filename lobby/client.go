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
	Name  string
	ID    string
	Conn  *websocket.Conn
	Lobby *Lobby
	Room  *Room
	Send  chan *Message
}

// NewClient creates a Client without a Room atached
func NewClient(lobby *Lobby, conn *websocket.Conn) *Client {
	return &Client{
		Name:  generators.Name(),
		ID:    uuid.NewString(),
		Lobby: lobby,
		Conn:  conn,
		Send:  make(chan *Message),
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
		msg := &Message{
			Content: &content,
		}
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			log.Println(err)
			return
		}
		switch msg.Type {
		case TypeLogin:
			var login Login
			if err := json.Unmarshal(content, &login); err != nil {
				log.Println(err)
			}
			if c.Room != nil {
				c.Room.Broadcast <- msg
			}
		case TypeChat:
			var chat Chat
			if err := json.Unmarshal(content, &chat); err != nil {
				log.Println(err)
			}
			c.Room.Broadcast <- msg
		case TypeJoin:
			var join Join
			if err := json.Unmarshal(content, &join); err != nil {
				log.Println(err)
			}
			c.Room.Broadcast <- msg
		case TypeLeave:
			var l Leave
			if err := json.Unmarshal(content, &l); err != nil {
				log.Println(err)
			}
			c.Room.Broadcast <- msg
		case TypeStartGame:
			c.Room.Game.StateChan <- StatusRunning
		default:
			log.Printf("unknown message type: %s", msg.Type)
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
