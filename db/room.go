package db

import (
	"errors"

	"github.com/paroar/battle-brush-backend/model"
)

//CREATE
func CreateRoom(r *model.Room) {
	rdb := roomRedisConnection()

	err := rdb.HSet(ctx, r.ID, r).Err()
	if err != nil {
		panic(err)
	}
}

//READ
func ReadRoom(id string) *model.Room {
	rdb := roomRedisConnection()

	res := rdb.HGetAll(ctx, id)
	if res.Err() != nil {
		panic(res.Err().Error())
	}

	var r model.Room
	if err := res.Scan(&r); err != nil {
		panic(err)
	}

	return &r
}

//UPDATE
func UpdateRoom(r *model.Room) {
	rdb := roomRedisConnection()

	err := rdb.HSet(ctx, r.ID, r).Err()
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

func AvailablePublicRoom() (string, error) {
	rdb := roomRedisConnection()

	publicRoomsID := publicRooms()
	for _, prID := range publicRoomsID {
		len, err := rdb.SCard(ctx, prID).Result()
		if err != nil {
			panic(err)
		}
		if len < 5 {
			return prID, nil
		}
	}
	return "", errors.New("Not available rooms")
}

func publicRooms() []string {
	rdb := roomRedisConnection()

	val, err := rdb.SMembers(ctx, "public").Result()
	if err != nil {
		panic(err)
	}

	return val
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
