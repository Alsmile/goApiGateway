package routers

import (
  "strconv"
  "github.com/kataras/iris"
  "github.com/alsmile/goApiGateway/utils"
  "github.com/alsmile/goApiGateway/controllers"
)

func SdkServer() {
  app := iris.New()
  app.Post("/api/site/apis/set", controllers.ApisSet)
  app.Post("/api/site/apis/delete", controllers.ApisDelete)

  strPort := strconv.Itoa(int(utils.GlobalConfig.Domain.SdkPort))
  app.Run(iris.Addr(":" + strPort))
}
