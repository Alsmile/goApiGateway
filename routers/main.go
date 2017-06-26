package routers

import (
  "strconv"
  "github.com/kataras/iris"
  "github.com/alsmile/goApiGateway/utils"
  "github.com/alsmile/goApiGateway/controllers"
  "github.com/kataras/iris/context"
)

func Start() {
  app := iris.New()

  app.OnErrorCode(iris.StatusNotFound, controllers.NotFound)

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
    admin.Get("/api/site/api/del", controllers.Auth, controllers.SiteApiDel)
    admin.Get("/api/site/api/list", controllers.Auth, controllers.SiteApiList)

    admin.Post("/api/site/api/list/by/domains", controllers.Auth, controllers.SiteApiListByDomains)

    admin.Get("/api/test", controllers.ProxyTest)
    admin.Post("/api/test", controllers.ProxyTest)
    admin.Put("/api/test", controllers.ProxyTest)
    admin.Delete("/api/test", controllers.ProxyTest)

    admin.Options("/api/{url:path}", func(ctx context.Context) {
      ctx.Header("Access-Control-Allow-Origin", "*")
      ctx.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, token, X-Requested-With")
      ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
      method := string(ctx.Method())
      if method == "OPTIONS" {
        ctx.JSON(true)
        return
      }
    })
  }

  app.Any("/{url:path}", controllers.ProxyDo)

  strPort := strconv.Itoa(int(utils.GlobalConfig.Domain.Port))
  app.Run(iris.Addr(":" + strPort))
}
