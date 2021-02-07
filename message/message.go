package message

// Message struct for websockets messages
type Message struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

// Common Message types
const (
	TypeLogin  = "Login"
	TypeLogout = "Logout"
	TypeJoin   = "Join"
	TypeLeave  = "Leave"
	TypeChat   = "Chat"
	TypeVote   = "Vote"
)
