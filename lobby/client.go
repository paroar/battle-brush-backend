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
		case TypeGameState:
			var gameState GameState
			if err := json.Unmarshal(content, &gameState); err != nil {
				log.Println(err)
			}
			if gameState.Command == StateStart {
				go c.Room.Game.StartGame()
			}
		case TypeChat:
			var chat Chat
			if err := json.Unmarshal(content, &chat); err != nil {
				log.Println(err)
			}
			c.Room.Broadcast <- msg
		case TypeImage:
			var img Image
			if err := json.Unmarshal(content, &img); err != nil {
				log.Println(err)
			}
			client, err := c.Room.getClient(img.UserID)
			if err != nil {
				log.Println(err)
			}
			drawing := &Drawing{
				Client: client,
				Img:    img.Img,
			}
			c.Room.Game.drawChan <- drawing
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
