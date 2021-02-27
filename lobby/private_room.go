package lobby

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// PrivateRoom struct
type PrivateRoom struct {
	ID      string
	Clients map[string]*Client
	Game    interface{}
}

// NewPrivateRoom creates a Room
func NewPrivateRoom() interface{} {
	clients := make(map[string]*Client)
	return &PrivateRoom{
		Clients: clients,
		ID:      uuid.NewString(),
		Game:    NewDrawGame(clients),
	}
}

//GetID returns the Room ID
func (r *PrivateRoom) GetID() string {
	return r.ID
}

// JoinClient joins the Client into the Room
func (r *PrivateRoom) JoinClient(c *Client) error {

	if !r.IsAvailable() {
		return errors.New("Room is full")
	}

	c.room = r
	r.Clients[c.id] = c
	r.BroadcastJoinLeave(c.name, c.id, "has joined")
	r.BroadcastClientNames()

	return nil
}

// LeaveClient leaves the Client from the Room
func (r *PrivateRoom) LeaveClient(id string) {
	if c, ok := r.Clients[id]; ok {
		delete(r.Clients, id)
		r.BroadcastJoinLeave(c.name, c.id, "has left")
		r.BroadcastClientNames()
	}
}

//Broadcast sends a message to all the Clients in the Room
func (r *PrivateRoom) Broadcast(msg *Message) {
	for clientid := range r.Clients {
		r.Clients[clientid].send <- msg
	}
}

//GetNames returns the names of all the Clients in the Room
func (r *PrivateRoom) GetNames() []string {
	names := []string{}
	for _, client := range r.Clients {
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

//GetGame returns game
func (r *PrivateRoom) GetGame() interface{} {
	return r.Game
}
