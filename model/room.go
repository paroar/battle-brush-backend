package model

import (
	"encoding/json"

	"github.com/google/uuid"
)

type Room struct {
	ID        string   `redis:"id"`
	PlayersID []string `redis:"playersid"`
	State     string   `redis:"state"`
	Type      string   `redis:"type"`
}

func NewRoom() *Room {
	return &Room{
		ID:        uuid.NewString(),
		PlayersID: make([]string, 0),
		State:     "Waiting",
	}
}

func UpdateRoom(id string, players []string, state string) *Room {
	return &Room{
		ID:        id,
		PlayersID: players,
		State:     state,
	}
}

func (r *Room) MarshalBinary() ([]byte, error) {
	return json.Marshal(r)
}
