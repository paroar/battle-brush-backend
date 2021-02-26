package lobby

//IRoom Room interface
type IRoom interface {
	GetID() string
	JoinClient(c *Client) error
	LeaveClient(c *Client)
	Broadcast(msg *Message)
	IsPublic() bool
	IsAvailable() bool
}

//Room Types
const (
	PrivateRoomType = iota
	PublicRoomType
)
