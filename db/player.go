package db

import (
	"log"

	"github.com/paroar/battle-brush-backend/model"
)

// CreatePlayer creates a player on redis database
func CreatePlayer(p *model.Player) {
	rdb := playerRedisConnection()

	err := rdb.HSet(ctx, p.ID, "player", p).Err()
	if err != nil {
		log.Println(err)
	}
}

// ReadPlayer reads a player from redis database
func ReadPlayer(id string) (*model.Player, error) {
	rdb := playerRedisConnection()

	res, err := rdb.HGetAll(ctx, id).Result()
	if err != nil {
		log.Println(err)
	}

	var p model.Player
	if err := p.UnmarshalBinary([]byte(res["player"])); err != nil {
		return nil, err
	}

	return &p, nil
}

// UpdatePlayer updates a player from redis database
func UpdatePlayer(p *model.Player) {
	rdb := playerRedisConnection()

	err := rdb.HSet(ctx, p.ID, "player", p).Err()
	if err != nil {
		log.Println(err)
	}

}

// DeletePlayer deletes a player from redis database and
// updates the room if it was in one
func DeletePlayer(id string) (*model.Player, error) {
	rdb := playerRedisConnection()

	player, err := ReadPlayer(id)
	if err != nil {
		log.Println(err)
	}

	roomid := player.RoomID
	if roomid != "" {
		room, err := ReadRoom(roomid)
		if err != nil {
			log.Println(err)
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
		log.Println(err)
	}

	return player, nil
}

// ReadPlayers gets players names of the room
func ReadPlayers(playersid []string) []*model.Player {
	players := []*model.Player{}
	for _, p := range playersid {
		player, err := ReadPlayer(p)
		if err != nil {
			log.Println(err)
		} else {
			players = append(players, player)
		}
	}

	return players
}
