package controllers
import (
  "gopkg.in/kataras/iris.v6"
)
func Index(ctx *iris.Context) {
  ctx.ServeFile("./admin/web/dist/index.html", true)
}

func Browser(ctx *iris.Context) {
  ctx.ServeFile("./admin/web/dist/browser.html", true)
}
