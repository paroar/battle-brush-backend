package games

//IGame interface
type IGame interface {
	StartGame()
}

// Common Game states
const (
	StateDrawing          = "Drawing"
	StateLoadingDrawing   = "LoadingDrawing"
	StateVoting           = "Voting"
	StateLoadingVoting    = "LoadingVoting"
	StateRecolecting      = "Recolecting"
	StateStart            = "Start"
	StateWaiting          = "Waiting"
	StateLoading          = "Loading"
	StateRecolectingVotes = "RecolectingVotes"
	StateLoadingWinner    = "LoadingWinner"
	StateWinner           = "Winner"
)
