package lobby

// GameState struct
type GameState struct {
	State   string `json:"gameState"`
	Command string `json:"command"`
}

//Do starts the game
func (g *GameState) Do(c *Client) {
	if g.Command == StateStart {
		go c.room.game.startGame()
	}
}
