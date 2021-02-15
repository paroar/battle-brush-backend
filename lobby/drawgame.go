package lobby

import (
	"log"
)

// DrawGame game
type DrawGame struct {
	players     map[*Client]bool
	theme       string
	state       string
	StateChan   chan string
	drawChan    chan string
	votingChan  chan int
	gameOptions Options
}

// Options struct for the game
type Options struct {
	DrawTime   int
	VotingTime int
	Rounds     int
}

var defaultGameOptions = &Options{
	DrawTime:   90,
	VotingTime: 20,
	Rounds:     3,
}

// NewDrawGame constructor
func NewDrawGame(players map[*Client]bool) *DrawGame {
	return &DrawGame{
		theme:       "",
		state:       StatusWaiting,
		StateChan:   make(chan string),
		drawChan:    make(chan string),
		votingChan:  make(chan int),
		gameOptions: *defaultGameOptions,
		players:     players,
	}
}

// RunGame gets the game going
func (d *DrawGame) Run() {
	for {
		select {
		case state := <-d.StateChan:
			d.changeState(state)
		case draw := <-d.drawChan:
			log.Println(draw)
		case vote := <-d.votingChan:
			log.Println(vote)
		}
	}
}

// ChangeState changes the state of the game
func (d *DrawGame) changeState(state string) {
	d.state = state
	msg := &Message{
		Type: TypeGameState,
		Content: GameState{
			State: state,
		},
	}
	d.broadcast(msg)
}

func (d *DrawGame) broadcast(msg *Message) {
	for player := range d.players {
		player.Send <- msg
	}
}
