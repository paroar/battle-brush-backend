package model

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/paroar/battle-brush-backend/generators"
)

// Player struct
type Player struct {
	ID     string `redis:"id" json:"id"`
	Name   string `redis:"name" json:"name"`
	RoomID string `redis:"roomid" json:"roomid"`
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
