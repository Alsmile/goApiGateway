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
  Website string `json:"website"`
  Cpu     int `json:"cpu"`
  Admin struct {
    Host string `json:"host"`
    Port uint16 `json:"port"`
  } `json:"admin"`
  User struct {
    LoginUrl  string `json:"loginUrl"`
    SignUpUrl string `json:"signUpUrl"`
    InfoUrl   string `json:"infoUrl"`
  } `json:"user"`
  Jwt    string `json:"jwt"`
  Secret string `json:"secret"`
  Mongo struct {
    Address        string `json:"address"`
    Database       string `json:"database"`
    User           string `json:"user"`
    Password       string `json:"password"`
    MaxConnections int `json:"maxConnections"`
    AcquireTimeout time.Duration `json:"acquireTimeout"`
    Mechanism      string `json:"mechanism"`
    Debug          bool `json:"debug"`
  } `json:"mongo"`
  PostgresSql struct {
    Host           string `json:"host"`
    Port           uint16 `json:"port"`
    Database       string `json:"database"`
    User           string `json:"user"`
    Password       string `json:"password"`
    MaxConnections int `json:"maxConnections"`
    AcquireTimeout time.Duration `json:"acquireTimeout"`
    ShowSql        bool `json:"showSql"`
    LogLevel       uint8 `json:"logLevel"`
    SslMode        string `json:"sslMode"`
  } `json:"postgresSql"`
  Redis struct {
    ConnectNum  int `json:"connectNum"`
    Address     string `json:"address"`
    Password    string `json:"password"`
    IdleTimeout int64 `json:"idleTimeout"`
    Db          string `json:"db"`
  } `json:"redis"`
  Log struct {
    Filename   string `json:"filename"`
    MaxSize    int `json:"maxSize"`
    MaxBackups int `json:"maxBackups"`
    MaxAge     int `json:"maxAge"`
  } `json:"log"`
  Email struct {
    Address  string `json:"address"`
    Port     int `json:"port"`
    User     string `json:"user"`
    Password string `json:"password"`
  } `json:"email"`
}

var GlobalConfig AppConfig

func ReadConfig() error {
  var defaultConfig AppConfig
  err := configor.Load(&defaultConfig, "./config/default.json")
  if err != nil {
    return err
  }

  files, err := WalkDir("./config", ".json")
  var config AppConfig
  for _, v := range files {
    err = configor.Load(&config, v)
    if err == nil && strings.Contains(v, "default.json") == false {
      jsonmp.PatchValue(defaultConfig, config, &GlobalConfig)
      defaultConfig = GlobalConfig
    }
  }
  // Second -> Duration
  GlobalConfig.PostgresSql.AcquireTimeout = GlobalConfig.PostgresSql.AcquireTimeout * time.Second

  return err
}
