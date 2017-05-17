package servers

import (
  "fmt"
  "strconv"
  "gopkg.in/kataras/iris.v6"
  "gopkg.in/kataras/iris.v6/adaptors/httprouter"
  "github.com/alsmile/goMicroServer/servers/controllers"
  "github.com/alsmile/goMicroServer/utils"
)


func Start() {
  app := iris.New()
  app.Adapt(httprouter.New())
  app.Any("/:key/*url", controllers.ProxyDo)
  fmt.Printf("[log]Proxy listen: %s:%d\r\n",utils.GlobalConfig.Gateway.Host, utils.GlobalConfig.Gateway.Port)
  strPort := strconv.Itoa(int(utils.GlobalConfig.Gateway.Port))
  app.Listen(utils.GlobalConfig.Gateway.Host + ":" + strPort)
}
