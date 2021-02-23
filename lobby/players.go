package lobby

// Players struct
type Players struct {
	UserNames []string `json:"usernames"`
}

//Do nothing
func (p *Players) Do(c *Client) {}
