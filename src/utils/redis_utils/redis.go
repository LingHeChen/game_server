package redis_utils

import (
	"game_server/src/frame"
	"fmt"

	"github.com/go-redis/redis/v8"
)


var RedisClinet *redis.Client



func init()  {
  redisAddr := fmt.Sprintf(
    "%s:%d",
    frame.Config.Redis.Host,
    frame.Config.Redis.Port,
  )
  RedisClinet = redis.NewClient(&redis.Options{
    Addr: redisAddr,
    Password: "",
    DB: frame.Config.Redis.DB,
  })
}
