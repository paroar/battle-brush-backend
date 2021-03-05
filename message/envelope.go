package message

// Envelope struct for websockets messages
type Envelope struct {
	Type    int         `json:"type"`
	Content interface{} `json:"content"`
}

// Content types
const (
	TypeLogin = iota
	TypeJoinLeave
	TypeChat
	TypePlayers
	TypeGameState
	TypeImage
	TypeTheme
	TypeVote
	TypeWinner
	TypeConnection
)
