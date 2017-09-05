package main

import (
	"log"
	"runtime"

	"github.com/alsmile/goApiGateway/db"
	"github.com/alsmile/goApiGateway/db/mongo"
	"github.com/alsmile/goApiGateway/routers"
	"github.com/alsmile/goApiGateway/utils"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	// 处理panic
	defer func() {
		if err := recover(); err != nil {
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
	if utils.GlobalConfig.Log.Filename != "" {
		log.SetOutput(&lumberjack.Logger{
			Filename:   utils.GlobalConfig.Log.Filename,
			MaxSize:    utils.GlobalConfig.Log.MaxSize, // mb
			MaxBackups: utils.GlobalConfig.Log.MaxBackups,
			MaxAge:     utils.GlobalConfig.Log.MaxAge, // days
		})
	}

	// cpu
	runtime.GOMAXPROCS(utils.GlobalConfig.CPU)

	// 连接数据库
	err = db.Init()
	if err != nil {
		log.Printf("[error]Db error: %v\r\n", err)
		return
	}
	defer mongo.MgoSession.Close()

	// 内部sdk服务
	go routers.SdkServer()

	// 后台管理web + proxy
	routers.Start()
}
