package redis

import (
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func Init() {
    RedisClient = redis.NewClient(&redis.Options{
        Addr:	  ":6379",
        Password: "", // no password set
        DB:		  0,  // use default DB
    })
}