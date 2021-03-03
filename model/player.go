package model

import (
	"github.com/google/uuid"
	"github.com/paroar/battle-brush-backend/generators"
)

type Player struct {
	ID   string `redis:"id"`
	Name string `redis:"name"`
}

func NewPlayer() *Player {
	return &Player{
		ID:   uuid.NewString(),
		Name: generators.Name(),
	}
}
