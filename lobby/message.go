package lobby

// Message struct for websockets messages
type Message struct {
	Type    string      `json:"type"`
	Content interface{} `json:"content"`
}

// Login struct
type Login struct {
	UserName string `json:"username"`
	ID       string `json:"userid"`
}

// JoinLeave struct
type JoinLeave struct {
	UserName string `json:"username"`
	ID       string `json:"userid"`
	Msg      string `json:"msg"`
}

// Chat struct
type Chat struct {
	Name string `json:"username"`
	Msg  string `json:"msg"`
}

// Players struct
type Players struct {
	UserNames []string `json:"usernames"`
}

// GameState struct
type GameState struct {
	State   string `json:"gameState"`
	Command string `json:"command"`
}

// Image struct
type Image struct {
	UserID   string `json:"userid"`
	UserName string `json:"username"`
	Img      string `json:"img"`
}

// Theme struct
type Theme struct {
	Theme string `json:"theme"`
}

// Vote struct
type Vote struct {
	Vote   float64 `json:"vote"`
	UserID string  `json:"userid"`
}

// Common Message types
const (
	TypeLogin     = "Login"
	TypeJoinLeave = "JoinLeave"
	TypeChat      = "Chat"
	TypePlayers   = "Players"
	TypeGameState = "GameState"
	TypeImage     = "Image"
	TypeTheme     = "Theme"
	TypeVote      = "Vote"
	TypeWinner    = "Winner"
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
