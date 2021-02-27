package lobby

//IRoom Room interface
type IRoom interface {
	GetID() string
	JoinClient(c *Client) error
	LeaveClient(id string)
	Broadcast(msg *Message)
	IsPublic() bool
	IsAvailable() bool
	GetGame() interface{}
}

//Room Types
const (
	PrivateRoomType = iota
	PublicRoomType
)
