package game

type DrawGame struct {
	ID      string
	Players []string
}

func NewDrawGame(id string, players []string) *DrawGame {
	return &DrawGame{
		ID:      id,
		Players: players,
	}
}

func (d *DrawGame) StartGame() {}
