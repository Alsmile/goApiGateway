package db

import (
  "github.com/alsmile/goMicroServer/db/mongo"
  "github.com/alsmile/goMicroServer/db/redis"
)

func Init() error {
  err := mongo.InitSession()
  if err != nil {
    return err
  }
  err = redis.NewPool()
  return err
}
