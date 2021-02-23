package lobby

// JoinLeave struct
type JoinLeave struct {
	UserName string `json:"username"`
	ID       string `json:"userid"`
	Msg      string `json:"msg"`
}

//Do nothing
func (jl *JoinLeave) Do(c *Client) {}
