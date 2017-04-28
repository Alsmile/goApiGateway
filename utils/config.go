package utils

import (
  "github.com/jinzhu/configor"
  "github.com/martingallagher/go-jsonmp"
  "strings"
  "time"
)

type AppConfig struct {
  Name    string `json:"name"`
  Version string `json:"version"`
  Cpu     int `json:"cpu"`
  Admin   struct {
            Host string `json:"host"`
            Port uint16 `json:"port"`
          } `json:"admin"`
  Db      struct {
             Host string `json:"host"`
             Port uint16 `json:"port"`
             Database string `json:"database"`
             User string `json:"user"`
             Password string `json:"password"`
             MaxConnections int `json:"maxConnections"`
             AcquireTimeout time.Duration `json:"acquireTimeout"`
           } `json:"db"`
  Log      struct{
             Filename string `json:"filename"`
             MaxSize int `json:"maxSize"`
             MaxBackups int `json:"maxBackups"`
             MaxAge int `json:"maxAge"`
           } `json:"log"`
}

func GetAppConfig() (AppConfig, error) {
  var appConfig AppConfig

  var defaultConfig AppConfig
  err := configor.Load(&defaultConfig, "./config/default.json")
  if err != nil {
    return appConfig, err
  }

  files, err := WalkDir("./config", ".json")
  var config AppConfig
  for _, v := range files {
    err = configor.Load(&config, v)
    if err == nil && strings.Contains(v, "default.json")  == false {
      jsonmp.PatchValue(defaultConfig, config, &appConfig)
      defaultConfig = appConfig
    }
  }
  // Second -> Duration
  appConfig.Db.AcquireTimeout = appConfig.Db.AcquireTimeout * time.Second

  return appConfig, err
}

