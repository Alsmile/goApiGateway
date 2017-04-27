package pq

import (
  "github.com/jackc/pgx"
  "github.com/alsmile/goMicroServer/utils"
)

var ConnPool *pgx.ConnPool

func GetConnPool() error {
  appConfig, err := utils.GetAppConfig()
  if err != nil {
    return err
  }

  config := pgx.ConnPoolConfig {
    ConnConfig: pgx.ConnConfig{
      Host: appConfig.Db.Host,
      Port: appConfig.Db.Port,
      Database: appConfig.Db.Database,
      User: appConfig.Db.User,
      Password: appConfig.Db.Password,
    },
    MaxConnections: appConfig.Db.MaxConnections,
    AcquireTimeout: appConfig.Db.AcquireTimeout,
  }

  ConnPool, err = pgx.NewConnPool(config)
  return err
}
