package db

import (
	"errors"

	"github.com/paroar/battle-brush-backend/model"
)

//CreateRoom creates a room on redis database
func CreateRoom(r *model.Room) {
	rdb := roomRedisConnection()

	err := rdb.HSet(ctx, r.ID, "room", r).Err()
	if err != nil {
		panic(err)
	}
}

//ReadRoom gets a room from redis database
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

//UpdateRoom updates a room from redis database
func UpdateRoom(r *model.Room) {
	rdb := roomRedisConnection()

	err := rdb.HSet(ctx, r.ID, "room", r).Err()
	if err != nil {
		panic(err)
	}
}

//DeleteRoom updates a room from redis database
func DeleteRoom(id string) {
	rdb := roomRedisConnection()

	err := rdb.Del(ctx, id).Err()
	if err != nil {
		panic(err)
	}
}

//AvailablePublicRoom gets the first available room if it exists from redis database
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
