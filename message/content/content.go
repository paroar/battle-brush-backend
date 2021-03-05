package content

// Connection struct
type Connection struct {
	Status   string `json:"status"`
	RoomID   string `json:"roomid"`
	RoomType string `json:"roomtype"`
}

// JoinLeave struct
type JoinLeave struct {
	UserName string `json:"username"`
	ID       string `json:"userid"`
	Msg      string `json:"msg"`
}

// Login struct
type Login struct {
	UserName string `json:"username"`
	ID       string `json:"userid"`
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
