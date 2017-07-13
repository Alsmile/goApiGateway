package utils

import (
  "strings"
  "os"
  "github.com/jinzhu/configor"
  "github.com/martingallagher/go-jsonmp"
)

type AppConfig struct {
  Name    string `json:"name"`
  Version string `json:"version"`
  Website string `json:"website"`
  Cpu     int `json:"cpu"`
  Domain struct {
    Domain      string `json:"domain"`
    AdminDomain string `json:"adminDomain"`
    Port        uint16 `json:"port"`
    SdkPort     uint16 `json:"sdkPort"`
  } `json:"domain"`
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
    Timeout        int `json:"timeout"`
    Mechanism      string `json:"mechanism"`
    Debug          bool `json:"debug"`
  } `json:"mongo"`
  Redis struct {
    ConnectNum  int `json:"connectNum"`
    Address     string `json:"address"`
    Password    string `json:"password"`
    Timeout     int `json:"timeout"`
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

  getEnvConfig()

  return err
}

func getEnvConfig()  {
  text := os.Getenv("ApiGateway_Website")
  if text != "" {
    GlobalConfig.Website = text
  }

  num := Int(os.Getenv("ApiGateway_Cpu"))
  if num > 0 {
    GlobalConfig.Cpu = num
  }

  text = os.Getenv("ApiGateway_Domain_Domain")
  if text != "" {
    GlobalConfig.Domain.Domain = text
  }

  text = os.Getenv("ApiGateway_Domain_AdminDomain")
  if text != "" {
    GlobalConfig.Domain.AdminDomain = text
  }

  num = Int(os.Getenv("ApiGateway_Domain_Port"))
  if num > 0 {
    GlobalConfig.Domain.Port = uint16(num)
  }

  num = Int(os.Getenv("ApiGateway_Domain_SdkPort"))
  if num > 0 {
    GlobalConfig.Domain.SdkPort = uint16(num)
  }

  text = os.Getenv("ApiGateway_User_LoginUrl")
  if text != "" {
    GlobalConfig.User.LoginUrl = text
  }

  text = os.Getenv("ApiGateway_User_SignUpUrl")
  if text != "" {
    GlobalConfig.User.SignUpUrl = text
  }

  text = os.Getenv("ApiGateway_User_InfoUrl")
  if text != "" {
    GlobalConfig.User.InfoUrl = text
  }

  text = os.Getenv("ApiGateway_Jwt")
  if text != "" {
    GlobalConfig.Jwt = text
  }

  text = os.Getenv("ApiGateway_Secret")
  if text != "" {
    GlobalConfig.Secret = text
  }

  text = os.Getenv("ApiGateway_Mongo_Address")
  if text != "" {
    GlobalConfig.Mongo.Address = text
  }

  text = os.Getenv("ApiGateway_Mongo_Database")
  if text != "" {
    GlobalConfig.Mongo.Database = text
  }

  text = os.Getenv("ApiGateway_Mongo_User")
  if text != "" {
    GlobalConfig.Mongo.User = text
  }

  text = os.Getenv("ApiGateway_Mongo_Password")
  if text != "" {
    GlobalConfig.Mongo.Password = text
  }

  num = Int(os.Getenv("ApiGateway_Mongo_MaxConnections"))
  if num > 0 {
    GlobalConfig.Mongo.MaxConnections = num
  }

  text = os.Getenv("ApiGateway_Mongo_Mechanism")
  if text != "" {
    GlobalConfig.Mongo.Mechanism = text
  }

  num = Int(os.Getenv("ApiGateway_Redis_ConnectNum"))
  if num > 0 {
    GlobalConfig.Redis.ConnectNum = num
  }

  text = os.Getenv("ApiGateway_Redis_Address")
  if text != "" {
    GlobalConfig.Redis.Address = text
  }

  text = os.Getenv("ApiGateway_Redis_Password")
  if text != "" {
    GlobalConfig.Redis.Password = text
  }

  text = os.Getenv("ApiGateway_Redis_Db")
  if text != "" {
    GlobalConfig.Redis.Db = text
  }

  text = os.Getenv("ApiGateway_Log_Filename")
  if text != "" {
    GlobalConfig.Log.Filename = text
  }

  num = Int(os.Getenv("ApiGateway_Log_MaxAge"))
  if num > 0 {
    GlobalConfig.Log.MaxAge = num
  }

  num = Int(os.Getenv("ApiGateway_Log_MaxSize"))
  if num > 0 {
    GlobalConfig.Log.MaxSize = num
  }

  num = Int(os.Getenv("ApiGateway_Log_MaxBackups"))
  if num > 0 {
    GlobalConfig.Log.MaxBackups = num
  }


  text = os.Getenv("ApiGateway_Email_Address")
  if text != "" {
    GlobalConfig.Email.Address = text
  }

  num = Int(os.Getenv("ApiGateway_Email_Port"))
  if num > 0 {
    GlobalConfig.Email.Port = num
  }

  text = os.Getenv("ApiGateway_Email_User")
  if text != "" {
    GlobalConfig.Email.User = text
  }

  text = os.Getenv("ApiGateway_Email_Password")
  if text != "" {
    GlobalConfig.Email.Password = text
  }
}
