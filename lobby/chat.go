package lobby

// Chat struct
type Chat struct {
	Name string `json:"username"`
	Msg  string `json:"msg"`
}

// Do broadcasts the msg
func (c *Chat) Do(client *Client) {
	msg := &Message{
		Type:    TypeChat,
		Content: c,
	}
	client.room.broadcast <- msg
}
