package utils

import (
	"reflect"
	"strings"

	"github.com/InVisionApp/conjungo"
	"github.com/jinzhu/configor"
)

// AppConfig 全局配置数据结构
type AppConfig struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Website string `json:"website"`
	CPU     int    `json:"cpu"`
	Domain  struct {
		Domain      string `json:"domain"`
		AdminDomain string `json:"adminDomain"`
		Port        uint16 `json:"port"`
		SdkPort     uint16 `json:"sdkPort"`
	} `json:"domain"`
	User struct {
		LoginURL    string `json:"loginUrl"`
		ProfileURL  string `json:"profileUrl"`
		ProfileType int    `json:"profileType"`
	} `json:"user"`
	Jwt    string `json:"jwt"`
	Secret string `json:"secret"`
	Mongo  struct {
		Address        string `json:"address"`
		Database       string `json:"database"`
		User           string `json:"user"`
		Password       string `json:"password"`
		MaxConnections int    `json:"maxConnections"`
		Timeout        int    `json:"timeout"`
		Mechanism      string `json:"mechanism"`
		Debug          bool   `json:"debug"`
	} `json:"mongo"`
	Redis struct {
		ConnectNum int    `json:"connectNum"`
		Address    string `json:"address"`
		Password   string `json:"password"`
		Timeout    int    `json:"timeout"`
		Db         string `json:"db"`
	} `json:"redis"`
	Log struct {
		Filename   string `json:"filename"`
		MaxSize    int    `json:"maxSize"`
		MaxBackups int    `json:"maxBackups"`
		MaxAge     int    `json:"maxAge"`
	} `json:"log"`
	Email struct {
		Address  string `json:"address"`
		Port     int    `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
	} `json:"email"`
}

// GlobalConfig 全局配置文件实例
var GlobalConfig AppConfig

// ReadConfig 读取全局配置文件
func ReadConfig() error {
	err := configor.Load(&GlobalConfig, "/etc/goApiGateway.json")
	// 文件不存在时，err不报错
	if GlobalConfig.Mongo.Address != "" || GlobalConfig.Redis.Address != "" {
		return nil
	}

	err = configor.Load(&GlobalConfig, "./config/default.json")
	if err != nil {
		return err
	}

	var config AppConfig
	opts := conjungo.NewOptions()
	opts.MergeFuncs.SetTypeMergeFunc(
		reflect.TypeOf(""), // string
		func(t, s reflect.Value, o *conjungo.Options) (reflect.Value, error) {
			iT, _ := t.Interface().(string)
			iS, _ := s.Interface().(string)
			if iS != "" {
				return reflect.ValueOf(iS), nil
			}
			return reflect.ValueOf(iT), nil
		},
	)

	opts.MergeFuncs.SetTypeMergeFunc(
		reflect.TypeOf(0), // int
		func(t, s reflect.Value, o *conjungo.Options) (reflect.Value, error) {
			iT, _ := t.Interface().(int)
			iS, _ := s.Interface().(int)
			if iS != 0 {
				return reflect.ValueOf(iS), nil
			}
			return reflect.ValueOf(iT), nil
		},
	)
	opts.MergeFuncs.SetTypeMergeFunc(
		reflect.TypeOf(true), // int
		func(t, s reflect.Value, o *conjungo.Options) (reflect.Value, error) {
			iT, _ := t.Interface().(bool)
			iS, _ := s.Interface().(bool)
			if iS {
				return reflect.ValueOf(iS), nil
			}
			return reflect.ValueOf(iT), nil
		},
	)
	opts.MergeFuncs.SetTypeMergeFunc(
		reflect.TypeOf(uint16(0)), // int
		func(t, s reflect.Value, o *conjungo.Options) (reflect.Value, error) {
			iT, _ := t.Interface().(uint16)
			iS, _ := s.Interface().(uint16)
			if iS > 0 {
				return reflect.ValueOf(iS), nil
			}
			return reflect.ValueOf(iT), nil
		},
	)
	conjungo.Merge(&GlobalConfig, config, opts)

	files, err := WalkDir("./config", ".json")
	for _, v := range files {
		if strings.Contains(v, "default.json") == false {
			configor.Load(&config, v)
			conjungo.Merge(&GlobalConfig, config, opts)
		}
	}

	return err
}
