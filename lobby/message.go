package lobby

// Message struct for websockets messages
type Message struct {
	Type    int         `json:"type"`
	Content interface{} `json:"content"`
}

// // Login struct
// type Login struct {
// 	UserName string `json:"username"`
// 	ID       string `json:"userid"`
// }

// // JoinLeave struct
// type JoinLeave struct {
// 	UserName string `json:"username"`
// 	ID       string `json:"userid"`
// 	Msg      string `json:"msg"`
// }

// // Chat struct
// type Chat struct {
// 	Name string `json:"username"`
// 	Msg  string `json:"msg"`
// }

// // Players struct
// type Players struct {
// 	UserNames []string `json:"usernames"`
// }

// // GameState struct
// type GameState struct {
// 	State   string `json:"gameState"`
// 	Command string `json:"command"`
// }

// // Image struct
// type Image struct {
// 	UserID   string `json:"userid"`
// 	UserName string `json:"username"`
// 	Img      string `json:"img"`
// }

// // Theme struct
// type Theme struct {
// 	Theme string `json:"theme"`
// }

// // Vote struct
// type Vote struct {
// 	Vote   float64 `json:"vote"`
// 	UserID string  `json:"userid"`
// }

// // Connection struct
// type Connection struct {
// 	Status   string `json:"status"`
// 	RoomID   string `json:"roomid"`
// 	RoomType string `json:"roomtype"`
// }

// // RoomCommand struct
// type RoomCommand struct {
// 	Command string `json:"command"`
// 	RoomID  string `json:"roomid"`
// }

// Common Message types
// const (
// 	TypeLogin = iota
// 	TypeJoinLeave
// 	TypeChat
// 	TypePlayers
// 	TypeGameState
// 	TypeImage
// 	TypeTheme
// 	TypeVote
// 	TypeWinner
// 	TypeConnection
// 	TypeRoomCommand
// )

// // TypeHandler to automate the type of message
// var TypeHandler = map[int]func() interface{}{
// 	TypeLogin:       func() interface{} { return &Login{} },
// 	TypeJoinLeave:   func() interface{} { return &JoinLeave{} },
// 	TypeChat:        func() interface{} { return &Chat{} },
// 	TypePlayers:     func() interface{} { return &Players{} },
// 	TypeGameState:   func() interface{} { return &GameState{} },
// 	TypeImage:       func() interface{} { return &Image{} },
// 	TypeTheme:       func() interface{} { return &Theme{} },
// 	TypeVote:        func() interface{} { return &Vote{} },
// 	TypeWinner:      func() interface{} { return &Image{} },
// 	TypeConnection:  func() interface{} { return &Connection{} },
// 	TypeRoomCommand: func() interface{} { return &RoomCommand{} },
// }

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
