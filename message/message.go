package message

// IContent interface
type IContent interface {
	NewContent() IContent
}

// Message struct for websockets messages
type Message struct {
	Type    string      `json:"type"`
	Content interface{} `json:"content"`
}

type Login struct {
	UserName string `json:"userName"`
	ID       string `json:"userid"`
}

type Join struct {
	UserName string `json:"userName"`
	ID       string `json:"userid"`
}

type Leave struct {
	UserName string `json:"userName"`
	ID       string `json:"userid"`
}

type Chat struct {
	Name string `json:"userName"`
	Msg  string `json:"msg"`
}

type Players struct {
	UserNames []string `json:"userNames"`
}

// NewMessage returns
func NewMessage(messageType string, content IContent) *Message {
	return &Message{
		Type:    messageType,
		Content: content,
	}
}

// Common Message types
const (
	TypeLogin   = "Login"
	TypeLogout  = "Logout"
	TypeJoin    = "Join"
	TypeLeave   = "Leave"
	TypeChat    = "Chat"
	TypeVote    = "Vote"
	TypePlayers = "Players"
)
