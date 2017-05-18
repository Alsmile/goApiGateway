package admin

import (
  "fmt"
  "strconv"
  "gopkg.in/kataras/iris.v6"
  "gopkg.in/kataras/iris.v6/adaptors/httprouter"
  "github.com/alsmile/goMicroServer/utils"
  "github.com/alsmile/goMicroServer/admin/controllers"
  proxy "github.com/alsmile/goMicroServer/servers/controllers"
)

func Start() {
  app := iris.New()
  app.Adapt(httprouter.New())
  admin := app.Party(utils.GlobalConfig.Domain.AdminDomain)
  {
    admin.StaticWeb("/assets", "./admin/web/dist/assets")

    admin.Get("/", controllers.Index)
    admin.Get("/browser.html", controllers.Browser)
    admin.Get("/captcha", controllers.Captcha)

    admin.Post("/api/login", controllers.Login)
    admin.Post("/api/signup", controllers.SignUp)
    admin.Post("/api/sign/active", controllers.SignActive)
    admin.Post("/api/forget/password", controllers.ForgetPassword)
    admin.Post("/api/sign/new/password", controllers.NewPassword)
    admin.Get("/api/user/profile", controllers.UserProfile)

    admin.Get("/api/sign/config", controllers.GetSignConfig)

    admin.Get("/api/site/list", controllers.Auth, controllers.SiteList)
    admin.Get("/api/site/get", controllers.Auth, controllers.SiteGet)
    admin.Post("/api/site/save", controllers.Auth, controllers.SiteSave)
    admin.Post("/api/site/api/save", controllers.Auth, controllers.SiteApiSave)
    admin.Get("/api/site/api/get", controllers.Auth, controllers.SiteApiGet)
    admin.Get("/api/site/api/list", controllers.Auth, controllers.SiteApiList)

    admin.OnError(iris.StatusNotFound, controllers.Index)
  }

  app.Any("/:key/*url", proxy.ProxyDo)

  fmt.Printf("[log]Listen: %s:%d\r\n", utils.GlobalConfig.Domain.Domain, utils.GlobalConfig.Domain.Port)
  strPort := strconv.Itoa(int(utils.GlobalConfig.Domain.Port))
  app.Listen(utils.GlobalConfig.Domain.Domain + ":" + strPort)
}
