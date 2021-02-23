package lobby

// Login struct
type Login struct {
	UserName string `json:"username"`
	ID       string `json:"userid"`
}

//Do nothing
func (l *Login) Do(c *Client) {}
