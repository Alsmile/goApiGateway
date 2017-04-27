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
  config, err := utils.GetAppConfig()
  if err != nil {
    log.Printf("[error]Load app config error: %v\r\n.", err)
    return
  }
  log.Printf("[config] %v\r\n", config)

  // 设置log
  log.SetOutput(&lumberjack.Logger{
    Filename:   config.Log.Filename,
    MaxSize:    config.Log.MaxSize, // mb
    MaxBackups: config.Log.MaxBackups,
    MaxAge:     config.Log.MaxAge, // days
  })

  // cpu
  runtime.GOMAXPROCS(config.Cpu)

  // 连接数据库
  err = db.Init()
  if err != nil {
    log.Printf("[error]Db error: %v\r\n", err)
    return
  }

  defer pq.ConnPool.Close()

  admin.Start()
}

