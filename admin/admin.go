package admin

import (
  "fmt"
  "strconv"
  "gopkg.in/kataras/iris.v6"
  "gopkg.in/kataras/iris.v6/adaptors/httprouter"
  "github.com/alsmile/goMicroServer/utils"
  "github.com/alsmile/goMicroServer/admin/controllers"
)

func Start() {
  appConfig, err := utils.GetAppConfig()
  if err != nil {
    return
  }

  app := iris.New()
  app.Adapt(httprouter.New())
  app.StaticWeb("/assets", "./admin/web/dist/assets")

  app.Get("/", controllers.Index)
  app.Get("/browser.html", controllers.Browser)
  strPort := strconv.Itoa(int(appConfig.Admin.Port))

  fmt.Printf("[log]Admin listen: %s:%d\r\n", appConfig.Admin.Host, appConfig.Admin.Port)
  app.Listen(appConfig.Admin.Host + ":" + strPort)
}
