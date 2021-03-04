package model

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/paroar/battle-brush-backend/generators"
)

type Player struct {
	ID     string `redis:"id"`
	Name   string `redis:"name"`
	RoomID string `redis:"roomid"`
}

func NewPlayer() *Player {
	return &Player{
		ID:   uuid.NewString(),
		Name: generators.Name(),
	}
}

func (p *Player) MarshalBinary() (data []byte, err error) {
	return json.Marshal(p)
}

func (p *Player) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &p); err != nil {
		return err
	}

	return nil
}
