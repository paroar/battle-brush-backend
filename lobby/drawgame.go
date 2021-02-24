package lobby

import (
	"errors"
	"log"
	"time"

	"github.com/paroar/battle-brush-backend/drawing"
	"github.com/paroar/battle-brush-backend/generators"
	"github.com/paroar/battle-brush-backend/utils"
)

// DrawGame game
type DrawGame struct {
	players     map[*Client]bool
	drawings    map[string]string
	votes       map[*Client][]float64
	theme       string
	state       string
	stateChan   chan string
	drawChan    chan *drawing.Drawing
	votingChan  chan *Vote
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
	VotingTime: 10,
}

// NewDrawGame constructor
func NewDrawGame(players map[*Client]bool) *DrawGame {
	return &DrawGame{
		state:       StateWaiting,
		drawings:    make(map[string]string),
		votes:       make(map[*Client][]float64),
		stateChan:   make(chan string),
		drawChan:    make(chan *drawing.Drawing, 10),
		votingChan:  make(chan *Vote),
		gameOptions: *defaultGameOptions,
		players:     players,
	}
}

// run gets the game going
func (d *DrawGame) run() {
	for {
		select {
		case state := <-d.stateChan:
			d.changeState(state)
		case draw := <-d.drawChan:
			d.addDrawing(draw)
		case vote := <-d.votingChan:
			d.addVote(vote)
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
	d.setRandomTheme()
	d.startDrawing()
	d.recolectDrawings()
	d.votingDrawings()
	d.winner()
	d.changeState(StateWaiting)
}

func (d *DrawGame) setRandomTheme() {
	d.theme = generators.Theme()
	theme := &Message{
		Type: TypeTheme,
		Content: Theme{
			Theme: d.theme,
		},
	}
	d.broadcast(theme)
}

func (d *DrawGame) startDrawing() {
	d.drawings = make(map[string]string, len(d.players))
	d.changeState(StateDrawing)
	time.Sleep(time.Duration(d.gameOptions.DrawTime) * time.Second)
}

func (d *DrawGame) recolectDrawings() {
	d.changeState(StateRecolecting)
	time.Sleep(time.Second)
}

func (d *DrawGame) addDrawing(drawing *drawing.Drawing) {
	d.drawings[drawing.ClientID] = drawing.Img
}

func (d *DrawGame) addVote(vote *Vote) {
	client, err := d.getPlayer(vote.UserID)
	if err != nil {
		log.Println(err)
		return
	}
	d.votes[client] = append(d.votes[client], vote.Vote)
}

func (d *DrawGame) getPlayer(userid string) (*Client, error) {
	for player := range d.players {
		if userid == player.id {
			return player, nil
		}
	}
	return nil, errors.New("Player not found")
}

func (d *DrawGame) votingDrawings() {
	d.changeState(StateLoading)
	time.Sleep(time.Second)

	for userid, img := range d.drawings {
		msg := &Message{
			Type: TypeImage,
			Content: Image{
				UserID: userid,
				Img:    img,
			},
		}
		d.broadcast(msg)

		d.changeState(StateVoting)
		time.Sleep(time.Duration(defaultGameOptions.VotingTime) * time.Second)

		d.changeState(StateRecolectingVotes)
		time.Sleep(time.Second)

		d.changeState(StateLoading)
		time.Sleep(time.Second)
	}
}

func (d *DrawGame) winner() {
	var scores = make(map[*Client]float64)
	for player, votes := range d.votes {
		avg := utils.Average(votes)
		scores[player] = avg
	}

	winner := d.maxScore(scores)

	msg := &Message{
		Type: TypeWinner,
		Content: Image{
			Img:      d.drawings[winner.id],
			UserID:   winner.id,
			UserName: winner.name,
		},
	}
	d.broadcast(msg)

	d.changeState(StateWinner)
	time.Sleep(10 * time.Second)
}

func (d *DrawGame) maxScore(scores map[*Client]float64) *Client {
	var playerWinner *Client
	var maxScore = 0.0
	for player, score := range scores {
		if score > maxScore {
			maxScore = score
			playerWinner = player
		}
	}
	return playerWinner
}
