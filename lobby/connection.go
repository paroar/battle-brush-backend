package lobby

// Connection struct
type Connection struct {
	Status   string `json:"status"`
	RoomID   string `json:"roomid"`
	RoomType string `json:"roomtype"`
}

// Do nothing
func (conn *Connection) Do(c *Client) {}
