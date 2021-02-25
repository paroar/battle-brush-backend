package lobby

//IRoom Room interface
type IRoom interface {
	GetID() string
	JoinClient(c *Client)
	LeaveClient(c *Client)
	Broadcast()
	IsPublic() bool
	IsAvailable() bool
}

//Room Types
const (
	PrivateRoomType = iota
	PublicRoomType
)
