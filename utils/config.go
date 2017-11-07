package utils

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

// AppConfig 全局配置数据结构
type AppConfig struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Website string `json:"website"`
	CPU     int    `json:"cpu" yaml:"cpu"`
	Domain  struct {
		Domain      string `json:"domain"`
		AdminDomain string `json:"adminDomain" yaml:"adminDomain"`
		Port        uint16 `json:"port"`
		SdkPort     uint16 `json:"sdkPort" yaml:"sdkPort"`
	} `json:"domain"`
	User struct {
		LoginURL    string `json:"loginUrl" yaml:"loginUrl"`
		ProfileURL  string `json:"profileUrl" yaml:"profileUrl"`
		ProfileType int    `json:"profileType" yaml:"profileType"`
	} `json:"user"`
	Jwt   string `json:"jwt"`
	Mongo struct {
		Address        string `json:"address"`
		Database       string `json:"database"`
		User           string `json:"user"`
		Password       string `json:"password"`
		MaxConnections int    `json:"maxConnections" yaml:"maxConnections"`
		Timeout        int    `json:"timeout"`
		Mechanism      string `json:"mechanism"`
		Debug          bool   `json:"debug"`
	} `json:"mongo"`
	Redis struct {
		ConnectNum int    `json:"connectNum" yaml:"connectNum"`
		Address    string `json:"address"`
		Password   string `json:"password"`
		Timeout    int    `json:"timeout"`
		Db         string `json:"db"`
	} `json:"redis"`
	Log struct {
		Filename   string `json:"filename"`
		MaxSize    int    `json:"maxSize" yaml:"maxSize"`
		MaxBackups int    `json:"maxBackups" yaml:"maxBackups"`
		MaxAge     int    `json:"maxAge" yaml:"maxAge"`
	} `json:"log"`
	Email struct {
		Address  string `json:"address"`
		Port     int    `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
	} `json:"email"`
	Plugins struct {
		IP   bool `json:"ip"`
		Rate bool `json:"rate"`
		Log  bool `json:"log"`
	} `json:"plugins"`
}

// GlobalConfig 全局配置文件实例
var GlobalConfig AppConfig

// ReadConfig 读取全局配置文件
func ReadConfig() error {
	files, err := WalkDir("./config", ".yaml")

	configFiles := make([]string, 0, 30)
	if b, _ := Exists("/etc/goApiGateway.yaml"); b {
		configFiles = append(configFiles, "/etc/goApiGateway.yaml")
	} else {
		configFiles = append(configFiles, files...)
	}

	for _, c := range configFiles {
		data, err := ioutil.ReadFile(c)
		if err == nil {
			yaml.Unmarshal(data, &GlobalConfig)
		}
	}

	if GlobalConfig.Mongo.Address == "" || GlobalConfig.Redis.Address == "" {
		return errors.New("Error in configs.")
	}

	getEnvConfig()

	return err
}

func getEnvConfig() {
	text := os.Getenv("SUBDOMAIN")
	if text != "" {
		GlobalConfig.Domain.Domain = strings.Replace(GlobalConfig.Domain.Domain, "<domain>", text, -1)
		GlobalConfig.User.LoginURL = strings.Replace(GlobalConfig.User.LoginURL, "<domain>", text, -1)
		GlobalConfig.User.ProfileURL = strings.Replace(GlobalConfig.User.ProfileURL, "<domain>", text, -1)
	}

	num := Int(os.Getenv("API_GATEWAY_PORT"))
	if num > 0 {
		GlobalConfig.Domain.Port = uint16(num)
	}

	text = os.Getenv("JWT")
	if text != "" {
		GlobalConfig.Jwt = text
	}
}
