package games

//DrawGame struct
type DrawGame struct {
	ID      string
	Players []string
}

// NewDrawGame constructor
func NewDrawGame(id string, players []string) *DrawGame {
	return &DrawGame{
		ID:      id,
		Players: players,
	}
}

// StartGame starts the game
func (d *DrawGame) StartGame() {}
