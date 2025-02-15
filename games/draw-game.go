package games

import (
	"log"
	"sort"
	"time"

	"github.com/paroar/battle-brush-backend/db"
	"github.com/paroar/battle-brush-backend/generators"
	"github.com/paroar/battle-brush-backend/message/content"
	"github.com/paroar/battle-brush-backend/model"
	"github.com/paroar/battle-brush-backend/utils"
	"github.com/paroar/battle-brush-backend/websocket"
)

//DrawGame struct
type DrawGame struct {
	ID      string
	Players []string
}

type score struct {
	ID       string
	AvgScore float64
}

const (
	LOADING_TIME     = 5
	DRAWING_TIME     = 45
	WINNER_TIME      = 10
	VOTING_TIME      = 10
	RECOLECTING_TIME = 1
	END_TIME         = 2
)

// NewDrawGame constructor
func NewDrawGame(id string, players []string) *DrawGame {
	return &DrawGame{
		ID:      id,
		Players: players,
	}
}

// StartGame starts the game
func (d *DrawGame) StartGame(l *websocket.Lobby) {
	d.setRandomTheme(l)

	room, err := db.ReadRoom(d.ID)
	if err != nil {
		log.Println(err)
		return
	}

	//Drawing
	d.changeState(room, l, model.StateLoadingDrawing, LOADING_TIME)
	d.changeState(room, l, model.StateDrawing, DRAWING_TIME)

	//Recolecting Img
	d.changeState(room, l, model.StateRecolecting, RECOLECTING_TIME)

	//Voting
	d.changeState(room, l, model.StateLoadingVoting, LOADING_TIME)
	d.voting(room, l)

	//Winner
	d.win(room, l)
	d.changeState(room, l, model.StateLoadingWinner, LOADING_TIME)
	d.changeState(room, l, model.StateWinner, WINNER_TIME)

	//Clean db
	d.cleaning()

	//Waiting
	d.changeState(room, l, model.StateLoading, END_TIME)
	d.changeState(room, l, model.StateWaiting, 0)
}

func (d *DrawGame) broadcastState(l *websocket.Lobby, state string) {
	msg := content.NewGameState(state)
	l.Broadcast(d.Players, msg)
}

func (d *DrawGame) setRandomTheme(l *websocket.Lobby) {
	theme := generators.Theme()
	msg := content.NewTheme(theme)
	l.Broadcast(d.Players, msg)
}

func (d *DrawGame) broadcastWinner(l *websocket.Lobby, img, playerid, username string) {
	msg := content.NewWinner(img, playerid, username)
	l.Broadcast(d.Players, msg)
}

func (d *DrawGame) winnerID() []score {
	scores := []score{}

	for _, p := range d.Players {
		votes := db.ReadVotes(p)
		avg := utils.Average(votes)
		scores = append(scores, score{ID: p, AvgScore: avg})
	}

	sort.SliceStable(scores, func(i, j int) bool {
		return scores[i].AvgScore > scores[j].AvgScore
	})

	return scores
}

func (d *DrawGame) changeState(room *model.Room, l *websocket.Lobby, state string, sec time.Duration) {
	room.State = state
	db.UpdateRoom(room)
	d.broadcastState(l, room.State)
	time.Sleep(sec * time.Second)
}

func (d *DrawGame) cleaning() {
	for _, p := range d.Players {
		db.DeleteDrawing(p)
		db.DeleteVote(p)
	}
}

func (d *DrawGame) win(room *model.Room, l *websocket.Lobby) {
	scores := d.winnerID()

	for _, s := range scores {
		player, err := db.ReadPlayer(s.ID)
		if err != nil {
			log.Println(err)
			continue
		}
		drawing, err := db.ReadDrawing(s.ID)
		if err != nil {
			log.Println(err)
			continue
		}
		d.broadcastWinner(l, drawing.Img, player.ID, player.Name)
		return
	}
}

func (d *DrawGame) voting(room *model.Room, l *websocket.Lobby) {
	for _, p := range d.Players {
		img, err := db.ReadDrawing(p)
		if err != nil {
			log.Println(err)
			continue
		}
		msg := content.NewImage(img.Img, img.PlayerID)
		l.Broadcast(d.Players, msg)

		d.changeState(room, l, model.StateVoting, VOTING_TIME)
		d.changeState(room, l, model.StateRecolectingVotes, RECOLECTING_TIME)
		d.changeState(room, l, model.StateLoading, RECOLECTING_TIME)
	}
}
