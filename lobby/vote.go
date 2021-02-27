package lobby

// Vote struct
type Vote struct {
	Vote   float64 `json:"vote"`
	UserID string  `json:"userid"`
}

//Do retrieves votes
func (v *Vote) Do(c *Client) {
	game := c.room.GetGame().(IGame)
	game.Vote(v)
}
