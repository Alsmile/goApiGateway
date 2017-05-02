package main

import (
  "log"
  "runtime"
  "gopkg.in/natefinch/lumberjack.v2"
  "github.com/alsmile/goMicroServer/utils"
  "github.com/alsmile/goMicroServer/db"
  "github.com/alsmile/goMicroServer/db/pq"
  "github.com/alsmile/goMicroServer/admin"
)

func main() {
  // 处理panic
  defer func(){
    if err := recover();err!=nil{
      log.Printf("[panic] %v\r\n", err)
    }
  }()

  // 读取全局配置文件
  err := utils.ReadConfig()
  if err != nil {
    log.Printf("[error]Load app config error: %v\r\n.", err)
    return
  }
  log.Printf("[config] %v\r\n", utils.GlobalConfig)

  // 设置log
  log.SetOutput(&lumberjack.Logger{
    Filename:   utils.GlobalConfig.Log.Filename,
    MaxSize:    utils.GlobalConfig.Log.MaxSize, // mb
    MaxBackups: utils.GlobalConfig.Log.MaxBackups,
    MaxAge:     utils.GlobalConfig.Log.MaxAge, // days
  })

  // cpu
  runtime.GOMAXPROCS(utils.GlobalConfig.Cpu)

  // 连接数据库
  err = db.Init()
  if err != nil {
    log.Printf("[error]Db error: %v\r\n", err)
    return
  }

  defer pq.ConnPool.Close()

  // 后台管理web
  admin.Start()
}

