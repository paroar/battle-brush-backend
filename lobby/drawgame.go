package lobby

import (
	"log"
	"time"

	"github.com/paroar/battle-brush-backend/generators"
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
	DrawTime:   10,
	VotingTime: 5,
}

// Drawing struct
type Drawing struct {
	Client *Client
	Img    string
}

// newDrawGame constructor
func newDrawGame(players map[*Client]bool) *DrawGame {
	return &DrawGame{
		state:       StateWaiting,
		drawings:    make(map[*Client]string),
		StateChan:   make(chan string),
		drawChan:    make(chan *Drawing, 10),
		votingChan:  make(chan int),
		gameOptions: *defaultGameOptions,
		players:     players,
	}
}

// run gets the game going
func (d *DrawGame) run() {
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
		player.send <- msg
	}
}

// startGame starts the game process
func (d *DrawGame) startGame() {
	d.theme = generators.Theme()
	d.drawings = make(map[*Client]string, len(d.players))
	d.changeState(StateDrawing)
	time.Sleep(time.Duration(d.gameOptions.DrawTime) * time.Second)
	d.changeState(StateRecolecting)
	time.Sleep(time.Second)
	d.broadcastImages()
	d.changeState(StateWaiting)
}

func (d *DrawGame) addDrawing(drawing *Drawing) {
	d.drawings[drawing.Client] = drawing.Img
}

func (d *DrawGame) broadcastImages() {
	d.changeState(StateVoting)
	for client := range d.drawings {
		msg := &Message{
			Type: TypeImage,
			Content: Image{
				UserID: client.id,
				Img:    d.drawings[client],
			},
		}
		d.broadcast(msg)
		time.Sleep(time.Duration(defaultGameOptions.VotingTime) * time.Second)
	}
}
