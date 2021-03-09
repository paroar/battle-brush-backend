package content

import (
	"fmt"

	"github.com/paroar/battle-brush-backend/message"
	"github.com/paroar/battle-brush-backend/model"
)

// Connection struct
type Connection struct {
	Status   string `json:"status"`
	RoomID   string `json:"roomid"`
	RoomType string `json:"roomtype"`
}

// NewConnection message constructor
func NewConnection(roomid, status, roomtype string) *message.Envelope {
	return &message.Envelope{
		Type: message.TypeConnection,
		Content: Connection{
			RoomID:   roomid,
			Status:   status,
			RoomType: roomtype,
		},
	}
}

// JoinLeave struct
type JoinLeave struct {
	UserName string `json:"username"`
	ID       string `json:"userid"`
	Msg      string `json:"msg"`
}

// NewJoinLeave message constructor
func NewJoinLeave(player *model.Player, msg string) *message.Envelope {
	return &message.Envelope{
		Type: TypeJoinLeave,
		Content: JoinLeave{
			UserName: player.Name,
			ID:       player.ID,
			Msg:      fmt.Sprintf("%s %s", player.Name, msg),
		},
	}
}

// Login struct
type Login struct {
	UserName string `json:"username"`
	ID       string `json:"userid"`
}

// NewLogin message constructor
func NewLogin(player *model.Player) *message.Envelope {
	return &message.Envelope{
		Type: TypeLogin,
		Content: Login{
			UserName: player.Name,
			ID:       player.ID,
		},
	}
}

// Players struct
type Players struct {
	Data []*model.Player `json:"data"`
}

// NewPlayers message constructor
func NewPlayers(players []*model.Player) *message.Envelope {
	return &message.Envelope{
		Type: TypePlayers,
		Content: Players{
			Data: players,
		},
	}
}

// GameState struct
type GameState struct {
	State   string `json:"gameState"`
	Command string `json:"command"`
}

// NewGameState message constructor
func NewGameState(state string) *message.Envelope {
	return &message.Envelope{
		Type: TypeGameState,
		Content: GameState{
			State: state,
		},
	}
}

// Image struct
type Image struct {
	UserID string `json:"userid"`
	Img    string `json:"img"`
}

// NewImage message constructor
func NewImage(img, playerid string) *message.Envelope {
	return &message.Envelope{
		Type: TypeImage,
		Content: Image{
			Img:    img,
			UserID: playerid,
		},
	}
}

// Winner struct
type Winner struct {
	UserID   string `json:"userid"`
	UserName string `json:"username"`
	Img      string `json:"img"`
}

// NewWinner message constructor
func NewWinner(img, playerid, username string) *message.Envelope {
	return &message.Envelope{
		Type: TypeWinner,
		Content: Winner{
			Img:      img,
			UserID:   playerid,
			UserName: username,
		},
	}
}

// Theme struct
type Theme struct {
	Theme string `json:"theme"`
}

// NewTheme message constructor
func NewTheme(theme string) *message.Envelope {
	return &message.Envelope{
		Type: TypeTheme,
		Content: Theme{
			Theme: theme,
		},
	}
}

// Chat struct
type Chat struct {
	Roomid   string `json:"roomid"`
	Playerid string `json:"playerid"`
	Username string `json:"username"`
	Msg      string `json:"msg"`
}

// NewChat message constructor
func NewChat(roomid, playerid, username, msg string) *message.Envelope {
	return &message.Envelope{
		Type: TypeChat,
		Content: Chat{
			Roomid:   roomid,
			Playerid: playerid,
			Username: username,
			Msg:      msg,
		},
	}
}
