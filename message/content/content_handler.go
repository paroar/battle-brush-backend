package content

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
)

// TypeHandler to automate the type of message
var TypeHandler = map[int]func() interface{}{
	TypeLogin:      func() interface{} { return &Login{} },
	TypeJoinLeave:  func() interface{} { return &JoinLeave{} },
	TypePlayers:    func() interface{} { return &Players{} },
	TypeGameState:  func() interface{} { return &GameState{} },
	TypeImage:      func() interface{} { return &Image{} },
	TypeTheme:      func() interface{} { return &Theme{} },
	TypeVote:       func() interface{} { return &Vote{} },
	TypeWinner:     func() interface{} { return &Image{} },
	TypeConnection: func() interface{} { return &Connection{} },
}
