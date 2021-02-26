package lobby

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// PublicRoom struct
type PublicRoom struct {
	ID        string
	Clients   map[*Client]bool
	broadcast chan *Message
}

// NewPublicRoom creates a Room
func NewPublicRoom() interface{} {
	return &PublicRoom{
		Clients:   make(map[*Client]bool),
		ID:        uuid.NewString(),
		broadcast: make(chan *Message),
	}
}

//GetID returns the Room ID
func (pr *PublicRoom) GetID() string {
	return pr.ID
}

// JoinClient joins the Client into the Room
func (pr *PublicRoom) JoinClient(c *Client) error {

	if !pr.IsAvailable() {
		return errors.New("Room is full")
	}

	pr.Clients[c] = true
	pr.BroadcastJoinLeave(c.name, c.id, "has joined")
	pr.BroadcastClientNames()

	return nil
}

// LeaveClient leaves the Client from the Room
func (pr *PublicRoom) LeaveClient(c *Client) {
	if _, ok := pr.Clients[c]; ok {
		delete(pr.Clients, c)
		pr.BroadcastJoinLeave(c.name, c.id, "has left")
		pr.BroadcastClientNames()
	}
}

//Broadcast sends a message to all the Clients in the Room
func (pr *PublicRoom) Broadcast(msg *Message) {
	for client := range pr.Clients {
		client.send <- msg
	}
}

//GetNames returns the names of all the Clients in the Room
func (pr *PublicRoom) GetNames() []string {
	names := []string{}
	for client := range pr.Clients {
		names = append(names, client.name)
	}
	return names
}

// BroadcastClientNames broadcasts the names of the Clients in the Room
func (pr *PublicRoom) BroadcastClientNames() {
	msg := &Message{
		Type: TypePlayers,
		Content: Players{
			UserNames: pr.GetNames(),
		},
	}
	pr.Broadcast(msg)
}

// BroadcastJoinLeave broadcasts when a Client joins or leaves the Room
func (pr *PublicRoom) BroadcastJoinLeave(username, userid, notification string) {
	msg := &Message{
		Type: TypeJoinLeave,
		Content: JoinLeave{
			UserName: username,
			ID:       userid,
			Msg:      fmt.Sprintf("%s %s", username, notification),
		},
	}
	pr.Broadcast(msg)
}

// IsPublic returns if the Room is Public
func (pr *PublicRoom) IsPublic() bool {
	return true
}

// IsAvailable returns if the Room is available
func (pr *PublicRoom) IsAvailable() bool {
	return true
}
