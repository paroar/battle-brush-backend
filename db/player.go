package db

import (
	"errors"

	"github.com/paroar/battle-brush-backend/model"
)

func CreatePlayer(p *model.Player) {
	rdb := playerRedisConnection()

	err := rdb.Set(ctx, p.ID, p.Name, 0).Err()
	if err != nil {
		panic(err)
	}
}

func ReadPlayer(id string) (string, error) {
	rdb := playerRedisConnection()

	val, err := rdb.Get(ctx, id).Result()
	if err != nil {
		return "", errors.New("player not found")
	}
	return val, nil
}

func removePlayer(id string) {
	rdb := playerRedisConnection()

	err := rdb.Del(ctx, id).Err()
	if err != nil {
		panic(err)
	}
}
