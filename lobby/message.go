package lobby

// Message struct for websockets messages
type Message struct {
	Type    int         `json:"type"`
	Content interface{} `json:"content"`
}

// Room Types
const (
	RoomTypePrivate = "Private"
	RoomTypePublic  = "Public"
)

// Room Commands
const (
	RoomCommandCreate     = "Create"
	RoomCommandJoinCreate = "JoinCreate"
)

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
