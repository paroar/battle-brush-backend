package lobby

// IGame is the interface of all games
type IGame interface {
	Run()
}

// Game statuses
const (
	StatusWaiting     = "Waiting"
	StatusDrawing     = "Drawing"
	StatusRecolecting = "Recolecting"
	StatusRunning     = "Running"
	StatusVoting      = "Voting"
	StatusEnd         = "End"
)
