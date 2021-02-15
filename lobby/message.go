package lobby

// IContent interface
type IContent interface {
	NewContent() IContent
}

// Message struct for websockets messages
type Message struct {
	Type    string      `json:"type"`
	Content interface{} `json:"content"`
}

// Login struct
type Login struct {
	UserName string `json:"userName"`
	ID       string `json:"userid"`
}

// Join struct
type Join struct {
	UserName string `json:"userName"`
	ID       string `json:"userid"`
}

// Leave struct
type Leave struct {
	UserName string `json:"userName"`
	ID       string `json:"userid"`
}

// Chat struct
type Chat struct {
	Name string `json:"userName"`
	Msg  string `json:"msg"`
}

// Players struct
type Players struct {
	UserNames []string `json:"userNames"`
}

// GameState struct
type GameState struct {
	State string `json:"gameState"`
}

// Common Message types
const (
	TypeLogin     = "Login"
	TypeLogout    = "Logout"
	TypeJoin      = "Join"
	TypeLeave     = "Leave"
	TypeChat      = "Chat"
	TypeVote      = "Vote"
	TypePlayers   = "Players"
	TypeStartGame = "StartGame"
	TypeGameState = "GameState"
)
