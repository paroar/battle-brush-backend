package message

type Envelope struct {
	Type    int         `json:"type"`
	Content interface{} `json:"content"`
}

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

const (
	RoomTypePrivate = "Private"
	RoomTypePublic  = "Public"
)
