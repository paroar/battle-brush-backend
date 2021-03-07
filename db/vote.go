package db

import (
	"log"
	"strconv"

	"github.com/paroar/battle-brush-backend/model"
)

//CreateVote creates a vote on redis database
func CreateVote(v *model.Vote) {
	rdb := voteRedisConnection()

	err := rdb.LPush(ctx, v.PlayerID, v.Vote).Err()
	if err != nil {
		log.Println(err)
	}
}

//ReadVotes creates a vote on redis database
func ReadVotes(id string) []float64 {
	rdb := voteRedisConnection()

	votesStr, err := rdb.LRange(ctx, id, 0, -1).Result()
	if err != nil {
		log.Println(err)
	}
	votes := []float64{}
	for _, v := range votesStr {
		f, _ := strconv.ParseFloat(v, 64)
		votes = append(votes, f)
	}
	return votes
}

//DeleteVote deletes a vote from redis database
func DeleteVote(id string) {
	rdb := voteRedisConnection()

	err := rdb.Del(ctx, id).Err()
	if err != nil {
		log.Println(err)
	}
}
