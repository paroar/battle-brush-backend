package model

import (
	"encoding/json"

	"github.com/google/uuid"
)

// Room struct
type Room struct {
	ID        string   `redis:"id"`
	PlayersID []string `redis:"playersid"`
	State     string   `redis:"state"`
	RoomType  string   `redis:"roomtype"`
}

//Room Types
const (
	RoomTypePrivate = "Private"
	RoomTypePublic  = "Public"
)

// NewRoom constructor
func NewRoom(playerid, roomType string) *Room {
	players := []string{playerid}
	return &Room{
		ID:        uuid.NewString(),
		PlayersID: players,
		State:     "Waiting",
		RoomType:  roomType,
	}
}

// UpdateRoom updates the room players and state
func (r *Room) UpdateRoom(players []string, state string) {
	r.PlayersID = players
	r.State = state
}

// MarshalBinary marshaler
func (r *Room) MarshalBinary() (data []byte, err error) {
	return json.Marshal(r)
}

// UnmarshalBinary unmarshaler
func (r *Room) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &r); err != nil {
		return err
	}

	return nil
}
