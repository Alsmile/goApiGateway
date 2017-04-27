package db

import (
  "github.com/alsmile/goMicroServer/db/pq"
)

func Init() error {
  err := pq.GetConnPool()
  if err != nil {
    return err
  }

  return err
}
