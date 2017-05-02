package db

import (
  "github.com/alsmile/goMicroServer/db/pq"
  "github.com/alsmile/goMicroServer/db/redis"
)

func Init() error {
  err := pq.GetConnPool()
  if err != nil {
    return err
  }
  err = redis.NewPool()
  return err
}
