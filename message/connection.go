package message

// Connection struct
type Connection struct {
	Status   string `json:"status"`
	RoomID   string `json:"roomid"`
	RoomType string `json:"roomtype"`
}
