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
		var content json.RawMessage
		msg := &Message{
			Content: &content,
		}
		err := c.conn.ReadJSON(&msg)
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
				go c.room.game.startGame()
			}
		case TypeChat:
			var chat Chat
			if err := json.Unmarshal(content, &chat); err != nil {
				log.Println(err)
			}
			c.room.broadcast <- msg
		case TypeVote:
			var vote Vote
			if err := json.Unmarshal(content, &vote); err != nil {
				log.Println(err)
			}
			c.room.game.votingChan <- &vote
		case TypeImage:
			var img Image
			if err := json.Unmarshal(content, &img); err != nil {
				log.Println(err)
			}
			client, err := c.room.getClient(img.UserID)
			if err != nil {
				log.Println(err)
			}
			drawing := &Drawing{
				Client: client,
				Img:    img.Img,
			}
			c.room.game.drawChan <- drawing
		case TypeRoomCommand:
			var r RoomCommand
			if err := json.Unmarshal(content, &r); err != nil {
				log.Println(err)
			}
			command := r.Command
			switch command {
			case RoomCommandCreate:
				room := c.lobby.CreatePrivateRoom(c)
				msg := &Message{
					Type: TypeRoomCommand,
					Content: RoomCommand{
						Command: RoomCommandCreate,
						RoomID:  room.ID,
					},
				}
				c.send <- msg
				break
			case RoomCommandJoinCreate:
				room := c.lobby.CreateOrJoinPublicRoom(c)
				msg := &Message{
					Type: TypeRoomCommand,
					Content: RoomCommand{
						Command: RoomCommandJoinCreate,
						RoomID:  room.ID,
					},
				}
				c.send <- msg
				break
			default:
				log.Printf("Unknown command: %s", r.Command)
			}

		default:
			log.Printf("Unknown message type: %s", msg.Type)
		}
	}
}

func (c *Client) disconnect() {
	if c.room != nil {
		c.room.leaveClientChan <- c
	}
	c.lobby.leaveClientChan <- c
	c.conn.Close()
}
