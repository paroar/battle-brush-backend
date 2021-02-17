package router

// RoomIDJSON json structure of room id
type RoomIDJSON struct {
	ID string `json:"roomid"`
}

// ClientIDJSON json structure of client id
type ClientIDJSON struct {
	ID string `json:"userid"`
}

// RoomJSON json structure of room creation
type RoomJSON struct {
	ID         string `json:"userid"`
	NumPlayers int    `json:"numPlayers"`
	Time       int    `json:"time"`
	Rounds     int    `json:"rounds"`
}

// JoinRoomJSON json structure to join a Room
type JoinRoomJSON struct {
	RoomID   string `json:"roomid"`
	ClientID string `json:"userid"`
}
