package lobby

// import (
// 	"time"

// 	"github.com/paroar/battle-brush-backend/drawing"
// 	"github.com/paroar/battle-brush-backend/generators"
// 	"github.com/paroar/battle-brush-backend/utils"
// )

// // DrawGame game
// type DrawGame struct {
// 	clients     map[string]*Client
// 	drawings    map[string]string
// 	votes       map[string][]float64
// 	theme       string
// 	state       string
// 	stateChan   chan string
// 	drawChan    chan *drawing.Drawing
// 	votingChan  chan *Vote
// 	gameOptions Options
// }

// // Options struct for the game
// type Options struct {
// 	DrawTime   int
// 	VotingTime int
// 	Rounds     int
// }

// var defaultGameOptions = &Options{
// 	DrawTime:   10,
// 	VotingTime: 10,
// }

// // NewDrawGame constructor
// func NewDrawGame(clients map[string]*Client) interface{} {
// 	return &DrawGame{
// 		state:       StateWaiting,
// 		drawings:    make(map[string]string),
// 		votes:       make(map[string][]float64),
// 		stateChan:   make(chan string),
// 		drawChan:    make(chan *drawing.Drawing, 10),
// 		votingChan:  make(chan *Vote),
// 		gameOptions: *defaultGameOptions,
// 		clients:     clients,
// 	}
// }

// // run gets the game going
// func (d *DrawGame) run() {
// 	for {
// 		select {
// 		case state := <-d.stateChan:
// 			d.changeState(state)
// 		case draw := <-d.drawChan:
// 			d.addDrawing(draw)
// 		case vote := <-d.votingChan:
// 			d.addVote(vote)
// 		}
// 	}
// }

// // ChangeState changes the state of the game
// func (d *DrawGame) changeState(state string) {
// 	d.state = state
// 	msg := &Message{
// 		Type: TypeGameState,
// 		Content: GameState{
// 			State: state,
// 		},
// 	}
// 	d.broadcast(msg)
// }

// func (d *DrawGame) broadcast(msg *Message) {
// 	for _, client := range d.clients {
// 		client.send <- msg
// 	}
// }

// // StartGame starts the game process
// func (d *DrawGame) StartGame() {
// 	go d.run()
// 	d.setRandomTheme()
// 	d.startDrawing()
// 	d.recolectDrawings()
// 	d.votingDrawings()
// 	d.winner()
// 	d.changeState(StateWaiting)
// }

// func (d *DrawGame) setRandomTheme() {
// 	d.theme = generators.Theme()
// 	theme := &Message{
// 		Type: TypeTheme,
// 		Content: Theme{
// 			Theme: d.theme,
// 		},
// 	}
// 	d.broadcast(theme)
// }

// func (d *DrawGame) startDrawing() {
// 	d.drawings = make(map[string]string, len(d.clients))
// 	d.changeState(StateDrawing)
// 	time.Sleep(time.Duration(d.gameOptions.DrawTime) * time.Second)
// }

// func (d *DrawGame) recolectDrawings() {
// 	d.changeState(StateRecolecting)
// 	time.Sleep(time.Second)
// }

// func (d *DrawGame) addDrawing(drawing *drawing.Drawing) {
// 	d.drawings[drawing.ClientID] = drawing.Img
// }

// func (d *DrawGame) addVote(vote *Vote) {
// 	client := d.clients[vote.UserID]
// 	d.votes[client.id] = append(d.votes[client.id], vote.Vote)
// }

// func (d *DrawGame) votingDrawings() {
// 	d.changeState(StateLoading)
// 	time.Sleep(time.Second)

// 	for userid, img := range d.drawings {
// 		msg := &Message{
// 			Type: TypeImage,
// 			Content: Image{
// 				UserID: userid,
// 				Img:    img,
// 			},
// 		}
// 		d.broadcast(msg)

// 		d.changeState(StateVoting)
// 		time.Sleep(time.Duration(defaultGameOptions.VotingTime) * time.Second)

// 		d.changeState(StateRecolectingVotes)
// 		time.Sleep(time.Second)

// 		d.changeState(StateLoading)
// 		time.Sleep(time.Second)
// 	}
// }

// func (d *DrawGame) winner() {
// 	var scores = make(map[string]float64)
// 	for clientid, votes := range d.votes {
// 		avg := utils.Average(votes)
// 		scores[clientid] = avg
// 	}

// 	winner := d.maxScore(scores)

// 	client := d.clients[winner]

// 	msg := &Message{
// 		Type: TypeWinner,
// 		Content: Image{
// 			Img:      d.drawings[client.id],
// 			UserID:   client.id,
// 			UserName: client.name,
// 		},
// 	}
// 	d.broadcast(msg)

// 	d.changeState(StateWinner)
// 	time.Sleep(10 * time.Second)
// }

// func (d *DrawGame) maxScore(scores map[string]float64) string {
// 	var maxScore float64
// 	var winnerid string
// 	for player, score := range scores {
// 		if score > maxScore {
// 			maxScore = score
// 			winnerid = player
// 		}
// 	}
// 	return winnerid
// }

// //Vote interface method for Votes
// func (d *DrawGame) Vote(v *Vote) {
// 	d.votingChan <- v
// }

// //Drawing interface method for Drawings
// func (d *DrawGame) Drawing(dw *drawing.Drawing) {
// 	d.drawChan <- dw
// }
