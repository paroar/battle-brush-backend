package games

//IGame interface
type IGame interface {
	StartGame()
}

// Common Game states
const (
	StateDrawing          = "Drawing"
	StateVoting           = "Voting"
	StateRecolecting      = "Recolecting"
	StateStart            = "Start"
	StateWaiting          = "Waiting"
	StateLoading          = "Loading"
	StateRecolectingVotes = "RecolectingVotes"
	StateWinner           = "Winner"
)
