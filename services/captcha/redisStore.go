package captcha
import (
  "github.com/alsmile/goMicroServer/db/redis"
  goRedis "github.com/garyburd/redigo/redis"
  "github.com/alsmile/goMicroServer/session"
)

type captchaRedisStore struct {
}

func (this *captchaRedisStore) Set(id string, digits []byte) {
  redisConn := redis.RedisPool.Get()
  defer redisConn.Close()

  _, err := redisConn.Do("SETEX", "captcha." + id, session.SessionMaxAge, digits)
  if err != nil {
  }
}

func (this *captchaRedisStore) Get(id string, clear bool) (digits []byte) {
  redisConn := redis.RedisPool.Get()
  defer redisConn.Close()

  digits, _ = goRedis.Bytes(redisConn.Do("GET", "captcha." + id))
  return
}
