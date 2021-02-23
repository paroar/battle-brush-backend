package lobby

// Theme struct
type Theme struct {
	Theme string `json:"theme"`
}

//Do nothing
func (t *Theme) Do(c *Client) {}
