package db

import (
  "github.com/alsmile/goApiGateway/db/mongo"
  "github.com/alsmile/goApiGateway/db/redis"
  "log"
)

func Init() error {
  err := mongo.InitSession()
  if err != nil {
    log.Printf("[error]mongo connect error: %v\r\n", err)
    return err
  }
  err = redis.NewPool()
  if err != nil {
    log.Printf("[error]redis connect error: %v\r\n", err)
  }
  return err
}
