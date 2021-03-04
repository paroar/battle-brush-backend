package model

import (
	"encoding/json"

	"github.com/google/uuid"
)

type Room struct {
	ID        string   `redis:"id"`
	PlayersID []string `redis:"playersid"`
	State     string   `redis:"state"`
	RoomType  string   `redis:"roomtype"`
}

func NewRoom(playerid, roomType string) *Room {
	players := []string{playerid}
	return &Room{
		ID:        uuid.NewString(),
		PlayersID: players,
		State:     "Waiting",
		RoomType:  roomType,
	}
}

func (r *Room) UpdateRoom(players []string, state string) {
	r.PlayersID = players
	r.State = state
}

func (r *Room) MarshalBinary() (data []byte, err error) {
	return json.Marshal(r)
}

func (r *Room) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &r); err != nil {
		return err
	}

	return nil
}
