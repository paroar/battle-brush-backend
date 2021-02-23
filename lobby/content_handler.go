package lobby

// Common Message types
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
	TypeRoomCommand
)

// TypeHandler to automate the type of message
var TypeHandler = map[int]func() interface{}{
	TypeLogin:       func() interface{} { return &Login{} },
	TypeJoinLeave:   func() interface{} { return &JoinLeave{} },
	TypeChat:        func() interface{} { return &Chat{} },
	TypePlayers:     func() interface{} { return &Players{} },
	TypeGameState:   func() interface{} { return &GameState{} },
	TypeImage:       func() interface{} { return &Image{} },
	TypeTheme:       func() interface{} { return &Theme{} },
	TypeVote:        func() interface{} { return &Vote{} },
	TypeWinner:      func() interface{} { return &Image{} },
	TypeConnection:  func() interface{} { return &Connection{} },
	TypeRoomCommand: func() interface{} { return &RoomCommand{} },
}
