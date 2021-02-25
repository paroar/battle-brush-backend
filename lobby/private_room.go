package lobby

import (
	"fmt"

	"github.com/google/uuid"
)

// PrivateRoom struct
type PrivateRoom struct {
	ID        string
	Clients   map[*Client]bool
	broadcast chan *Message
}

// NewPrivateRoom creates a Room
func NewPrivateRoom() interface{} {
	return &PrivateRoom{
		Clients:   make(map[*Client]bool),
		ID:        uuid.NewString(),
		broadcast: make(chan *Message),
	}
}

//GetID returns the Room ID
func (r *PrivateRoom) GetID() string {
	return r.ID
}

// JoinClient joins the Client into the Room
func (r *PrivateRoom) JoinClient(c *Client) {
	r.Clients[c] = true
	r.BroadcastJoinLeave(c.name, c.id, "has joined")
	r.BroadcastClientNames()
}

// LeaveClient leaves the Client from the Room
func (r *PrivateRoom) LeaveClient(c *Client) {
	if _, ok := r.Clients[c]; ok {
		delete(r.Clients, c)
		r.BroadcastJoinLeave(c.name, c.id, "has left")
		r.BroadcastClientNames()
	}
}

//Broadcast sends a message to all the Clients in the Room
func (r *PrivateRoom) Broadcast(msg *Message) {
	for client := range r.Clients {
		client.send <- msg
	}
}

//GetNames returns the names of all the Clients in the Room
func (r *PrivateRoom) GetNames() []string {
	names := []string{}
	for client := range r.Clients {
		names = append(names, client.name)
	}
	return names
}

// BroadcastClientNames broadcasts the names of the Clients in the Room
func (r *PrivateRoom) BroadcastClientNames() {
	msg := &Message{
		Type: TypePlayers,
		Content: Players{
			UserNames: r.GetNames(),
		},
	}
	r.Broadcast(msg)
}

// BroadcastJoinLeave broadcasts when a Client joins or leaves the Room
func (r *PrivateRoom) BroadcastJoinLeave(username, userid, notification string) {
	msg := &Message{
		Type: TypeJoinLeave,
		Content: JoinLeave{
			UserName: username,
			ID:       userid,
			Msg:      fmt.Sprintf("%s %s", username, notification),
		},
	}
	r.Broadcast(msg)
}

// IsPublic returns if the Room is Public
func (r *PrivateRoom) IsPublic() bool {
	return false
}

// IsAvailable returns if the Room is available
func (r *PrivateRoom) IsAvailable() bool {
	return true
}
