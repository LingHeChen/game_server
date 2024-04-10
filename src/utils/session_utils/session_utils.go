package session_utils

import (
	"game_server/src/frame"
	"game_server/src/utils/redis_utils"
	"context"
	"time"
)


const (
  SESSION_DATA_KEY = "session:"
  DEFAULT_EXPIRE = 24 * time.Hour
)


func GetSessionData(sessionId string) (map[string]string, error) {
  ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
  defer cancel()
  key := SESSION_DATA_KEY + sessionId
  hashValue, err := redis_utils.RedisClinet.HGetAll(ctx, key).Result()
  if err != nil {
    return nil, err
  }
  return hashValue, nil
}


func SetSessionData(
  sessionId string,
  hashKey string,
  value string,
  timeout time.Duration,
  expire time.Duration,
) error {
  ctx := context.Background()
  var cancel context.CancelFunc
  if timeout > 0 {
    ctx, cancel = context.WithTimeout(ctx, timeout)
    defer cancel()
  }
  key := SESSION_DATA_KEY + sessionId
  if _, err := redis_utils.RedisClinet.HSet(ctx, key, hashKey, value).Result(); err != nil {
    frame.Logger.Error("Error Setting Session Data")
    return err
  }
  if expire < 1 {
    expire = DEFAULT_EXPIRE
  }
  if ok, err := redis_utils.RedisClinet.Expire(ctx, key, expire).Result(); !ok {
    frame.Logger.Error("Error Setting Sesssion Expiration")
    return err
  }
  return nil
}
