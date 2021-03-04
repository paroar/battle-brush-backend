package db

import (
	"errors"

	"github.com/paroar/battle-brush-backend/model"
)

//CREATE
func CreateRoom(r *model.Room) {
	rdb := roomRedisConnection()

	err := rdb.HSet(ctx, r.ID, "room", r).Err()
	if err != nil {
		panic(err)
	}
}

//READ
func ReadRoom(id string) (*model.Room, error) {
	rdb := roomRedisConnection()

	res, err := rdb.HGetAll(ctx, id).Result()
	if err != nil {
		panic(err)
	}

	var r model.Room
	if err := r.UnmarshalBinary([]byte(res["room"])); err != nil {
		return nil, err
	}

	return &r, nil
}

//UPDATE
func UpdateRoom(r *model.Room) {
	rdb := roomRedisConnection()

	err := rdb.HSet(ctx, r.ID, "room", r).Err()
	if err != nil {
		panic(err)
	}
}

//DELETE
func deleteRoom(id string) {
	rdb := roomRedisConnection()

	err := rdb.Del(ctx, id).Err()
	if err != nil {
		panic(err)
	}
}

func ReadRoomPlayers(id string) []string {
	rdb := roomRedisConnection()

	val, err := rdb.SMembers(ctx, id).Result()
	if err != nil {
		panic(err)
	}

	return val
}

func AvailablePublicRoom() (*model.Room, error) {
	rdb := roomRedisConnection()

	iterator := rdb.Scan(ctx, 0, "*", 1).Iterator()
	for iterator.Next(ctx) {

		roomid := iterator.Val()

		res, err := rdb.HGetAll(ctx, roomid).Result()
		if err != nil {
			panic(err)
		}

		var r model.Room
		if err := r.UnmarshalBinary([]byte(res["room"])); err != nil {
			return nil, err
		}

		if r.RoomType == "Public" && len(r.PlayersID) < 5 {
			return &r, nil
		}

	}

	return nil, errors.New("Not available rooms")
}

func AddPrivateRoom(id string) {
	rdb := roomRedisConnection()

	err := rdb.SAdd(ctx, "private", id).Err()
	if err != nil {
		panic(err)
	}
}

func RemovePrivateRoom(id string) {
	rdb := roomRedisConnection()

	err := rdb.SRem(ctx, "private", id).Err()
	if err != nil {
		panic(err)
	}
	deleteRoom(id)
}

func AddPublicRoom(id string) {
	rdb := roomRedisConnection()

	err := rdb.SAdd(ctx, "public", id).Err()
	if err != nil {
		panic(err)
	}
}

func RemovePublicRoom(id string) {
	rdb := roomRedisConnection()

	err := rdb.SRem(ctx, "public", id).Err()
	if err != nil {
		panic(err)
	}
	deleteRoom(id)
}
