package pq

import (
  "github.com/jackc/pgx"
  "github.com/alsmile/goMicroServer/utils"
)

var ConnPool *pgx.ConnPool

func GetConnPool() (err error) {
  config := pgx.ConnPoolConfig {
    ConnConfig: pgx.ConnConfig{
      Host: utils.GlobalConfig.PostgresSql.Host,
      Port: utils.GlobalConfig.PostgresSql.Port,
      Database: utils.GlobalConfig.PostgresSql.Database,
      User: utils.GlobalConfig.PostgresSql.User,
      Password: utils.GlobalConfig.PostgresSql.Password,
    },
    MaxConnections: utils.GlobalConfig.PostgresSql.MaxConnections,
    AcquireTimeout: utils.GlobalConfig.PostgresSql.AcquireTimeout,
  }
  ConnPool, err = pgx.NewConnPool(config)

  return
}
