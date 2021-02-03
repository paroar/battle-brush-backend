package router

// LobbyRoomsJSON json structure
type LobbyRoomsJSON struct {
	Rooms []LobbyRoomJSON `json:"rooms"`
}

// LobbyRoomJSON json structure
type LobbyRoomJSON struct {
	ID      string   `json:"id"`
	Clients []string `json:"clients"`
}

// RoomIDJSON json structure of room id
type RoomIDJSON struct {
	ID string `json:"roomid"`
}

// ClientIDJSON json structure of client id
type ClientIDJSON struct {
	ID string `json:"userid"`
}

// JoinRoomJSON json structure to join a Room
type JoinRoomJSON struct {
	RoomID   string `json:"roomid"`
	ClientID string `json:"userid"`
}

// ClientsIDJSON json structure of Clients ID
type ClientsIDJSON struct {
	ClientsID []string
}
