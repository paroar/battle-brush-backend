package lobby

import (
	"log"
	"time"
)

// DrawGame game
type DrawGame struct {
	players     map[*Client]bool
	drawings    map[*Client]string
	theme       string
	state       string
	StateChan   chan string
	drawChan    chan *Drawing
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
	DrawTime:   5,
	VotingTime: 5,
	Rounds:     1,
}

// Drawing struct
type Drawing struct {
	Client *Client
	Img    string
}

// NewDrawGame constructor
func NewDrawGame(players map[*Client]bool) *DrawGame {
	return &DrawGame{
		theme:       "",
		state:       StatusWaiting,
		drawings:    make(map[*Client]string),
		StateChan:   make(chan string),
		drawChan:    make(chan *Drawing, 10),
		votingChan:  make(chan int),
		gameOptions: *defaultGameOptions,
		players:     players,
	}
}

// Run gets the game going
func (d *DrawGame) Run() {
	for {
		select {
		case state := <-d.StateChan:
			d.changeState(state)
		case draw := <-d.drawChan:
			d.addDrawing(draw)
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

// StartGame starts the game process
func (d *DrawGame) StartGame() {
	d.drawings = map[*Client]string{}
	var rounds = d.gameOptions.Rounds
	for rounds > 0 {
		d.changeState(StatusDrawing)
		time.Sleep(time.Duration(d.gameOptions.DrawTime) * time.Second)
		d.changeState(StatusRecolecting)
		time.Sleep(time.Second)
		d.broadcastImages()
		rounds--
	}
	d.changeState(StatusWaiting)
}

func (d *DrawGame) addDrawing(drawing *Drawing) {
	d.drawings[drawing.Client] = drawing.Img
}

func (d *DrawGame) broadcastImages() {
	d.changeState(StatusVoting)
	for client := range d.drawings {
		msg := &Message{
			Type: TypeImage,
			Content: Image{
				UserID: client.ID,
				Img:    d.drawings[client],
			},
		}
		d.broadcast(msg)
		time.Sleep(2 * time.Second)
	}
}
