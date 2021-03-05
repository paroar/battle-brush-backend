package model

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/paroar/battle-brush-backend/generators"
)

// Player struct
type Player struct {
	ID     string `redis:"id"`
	Name   string `redis:"name"`
	RoomID string `redis:"roomid"`
}

// NewPlayer constructor
func NewPlayer() *Player {
	return &Player{
		ID:   uuid.NewString(),
		Name: generators.Name(),
	}
}

// MarshalBinary marshaler
func (p *Player) MarshalBinary() (data []byte, err error) {
	return json.Marshal(p)
}

// UnmarshalBinary unmarshaler
func (p *Player) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &p); err != nil {
		return err
	}

	return nil
}
