package lobby

import "log"

// RoomCommand struct
type RoomCommand struct {
	Command string `json:"command"`
	RoomID  string `json:"roomid"`
}

//Do command switch
func (rc *RoomCommand) Do(c *Client) {
	switch rc.Command {
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
		log.Printf("Unknown command: %s", rc.Command)
	}
}
