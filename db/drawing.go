package db

import (
	"log"

	"github.com/paroar/battle-brush-backend/model"
)

//CreateDrawing creates a drawing on redis database
func CreateDrawing(d *model.Drawing) {
	rdb := imgRedisConnection()

	err := rdb.HSet(ctx, d.PlayerID, "drawing", d).Err()
	if err != nil {
		log.Println(err)
	}
}

//ReadDrawing gets a drawing from redis database
func ReadDrawing(id string) (*model.Drawing, error) {
	rdb := imgRedisConnection()

	res, err := rdb.HGetAll(ctx, id).Result()
	if err != nil {
		log.Println(err)
	}

	var d model.Drawing
	if err := d.UnmarshalBinary([]byte(res["drawing"])); err != nil {
		return nil, err
	}

	return &d, nil
}

//DeleteDrawing deletes a vote from redis database
func DeleteDrawing(id string) {
	rdb := imgRedisConnection()

	err := rdb.Del(ctx, id).Err()
	if err != nil {
		log.Println(err)
	}
}
