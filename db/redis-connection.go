package db

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
)

var (
	REDIS_PORT      = os.Getenv("REDIS_PORT")
	REDIS_URL       = os.Getenv("REDIS_URL")
	REDIS_PASSWORD  = os.Getenv("REDIS_PASSWORD")
	REDIS_PLAYER_DB = os.Getenv("REDIS_PLAYER_DB")
	REDIS_ROOM_DB   = os.Getenv("REDIS_ROOM_DB")
	REDIS_IMG_DB    = os.Getenv("REDIS_IMG_DB")
	REDIS_VOTE_DB   = os.Getenv("REDIS_VOTE_DB")
)

var ctx = context.Background()

func playerRedisConnection() *redis.Client {
	redisDB, err := strconv.Atoi(REDIS_PLAYER_DB)
	if err != nil {
		panic(err)
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", REDIS_URL, REDIS_PORT),
		Password: REDIS_PASSWORD,
		DB:       redisDB,
	})
	return rdb
}

func roomRedisConnection() *redis.Client {
	redisDB, err := strconv.Atoi(REDIS_ROOM_DB)
	if err != nil {
		panic(err)
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", REDIS_URL, REDIS_PORT),
		Password: REDIS_PASSWORD,
		DB:       redisDB,
	})
	return rdb
}

func imgRedisConnection() *redis.Client {
	redisDB, err := strconv.Atoi(REDIS_IMG_DB)
	if err != nil {
		panic(err)
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", REDIS_URL, REDIS_PORT),
		Password: REDIS_PASSWORD,
		DB:       redisDB,
	})
	return rdb
}

func voteRedisConnection() *redis.Client {
	redisDB, err := strconv.Atoi(REDIS_VOTE_DB)
	if err != nil {
		panic(err)
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", REDIS_URL, REDIS_PORT),
		Password: REDIS_PASSWORD,
		DB:       redisDB,
	})
	return rdb
}
