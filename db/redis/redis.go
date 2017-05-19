package redis

import (
  "time"
  "github.com/garyburd/redigo/redis"
  "github.com/alsmile/goApiGateway/utils"
)

var RedisPool *redis.Pool

func NewPool() (err error) {
  RedisPool = &redis.Pool{
    MaxIdle:     utils.GlobalConfig.Redis.ConnectNum * 2,
    MaxActive:   utils.GlobalConfig.Redis.ConnectNum,
    IdleTimeout: time.Duration(utils.GlobalConfig.Redis.IdleTimeout) * time.Second,
    Dial: func() (redis.Conn, error) {
      var c redis.Conn
      c, err = redis.Dial("tcp", utils.GlobalConfig.Redis.Address)
      if err != nil {
        return nil, err
      }
      if _, err = c.Do("AUTH", utils.GlobalConfig.Redis.Password); err != nil {
        c.Close()
        return nil, err
      }
      // select db
      c.Do("SELECT", utils.GlobalConfig.Redis.Db)
      return c, nil
    },
  }

  // 调用Get()执行RedisPool.Dial连接redis 。确认连接有效。
  defer RedisPool.Get().Close()

  return
}
