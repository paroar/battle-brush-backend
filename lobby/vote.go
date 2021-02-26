package lobby

// Vote struct
type Vote struct {
	Vote   float64 `json:"vote"`
	UserID string  `json:"userid"`
}

//Do retrieves votes
func (v *Vote) Do(c *Client) {
	// c.room.game.votingChan <- v
}
