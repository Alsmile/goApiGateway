package db

import (
  "github.com/alsmile/goApiGateway/db/mongo"
  "github.com/alsmile/goApiGateway/db/redis"
)

func Init() error {
  err := mongo.InitSession()
  if err != nil {
    return err
  }
  err = redis.NewPool()
  return err
}
