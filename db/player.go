package db

import (
	"github.com/paroar/battle-brush-backend/model"
)

func CreatePlayer(p *model.Player) {
	rdb := playerRedisConnection()

	err := rdb.HSet(ctx, p.ID, "player", p).Err()
	if err != nil {
		panic(err)
	}
}

func ReadPlayer(id string) (*model.Player, error) {
	rdb := playerRedisConnection()

	res, err := rdb.HGetAll(ctx, id).Result()
	if err != nil {
		panic(err)
	}

	var p model.Player
	if err := p.UnmarshalBinary([]byte(res["player"])); err != nil {
		return nil, err
	}

	return &p, nil
}

func UpdatePlayer(p *model.Player) {
	rdb := playerRedisConnection()

	err := rdb.HSet(ctx, p.ID, "player", p).Err()
	if err != nil {
		panic(err)
	}

}

func DeletePlayer(id string) (*model.Player, error) {
	rdb := playerRedisConnection()

	player, err := ReadPlayer(id)
	if err != nil {
		panic(err)
	}

	roomid := player.RoomID
	if roomid != "" {
		room, err := ReadRoom(roomid)
		if err != nil {
			panic(err)
		}
		updatedPlayers := []string{}
		players := room.PlayersID
		for _, p := range players {
			if p != id {
				updatedPlayers = append(updatedPlayers, p)
			}
		}
		room.UpdateRoom(updatedPlayers, room.State)
		UpdateRoom(room)
	}

	err = rdb.Del(ctx, id).Err()
	if err != nil {
		panic(err)
	}

	return player, nil
}
