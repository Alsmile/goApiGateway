package controllers
import (
  "gopkg.in/kataras/iris.v6"
  "github.com/garyburd/redigo/redis"
  "github.com/dchest/captcha"
  "github.com/alsmile/goApiGateway/session"
  myCaptcha "github.com/alsmile/goApiGateway/services/captcha"
  "strings"
)


func Index(ctx *iris.Context) {
  ctx.SetStatusCode(iris.StatusOK)
  ctx.ServeFile("./admin/web/dist/index.html", true)
}

func Browser(ctx *iris.Context) {
  ctx.ServeFile("./admin/web/dist/browser.html", true)
}

func NotFound(ctx *iris.Context) {
  if strings.HasPrefix(ctx.Path(), "/api/") {
    ret := make(map[string]interface{})
    ret["error"] = "请求错误（Not found）：" + ctx.Path()
    ctx.SetStatusCode(iris.StatusNotFound)
    ServeJson(ctx,ret)
  } else {
    ctx.SetStatusCode(iris.StatusFound)
    Index(ctx)
  }
}

func Captcha(ctx *iris.Context) {
  captchaId, _ := redis.String(session.GetSession(ctx, myCaptcha.CaptchaSessionName))

  // Delete the old.
  if captchaId != "" {
    captcha.VerifyString(captchaId, "")
  }

  captchaId = captcha.New()
  session.SetSession(ctx, myCaptcha.CaptchaSessionName, captchaId)
  ctx.SetHeader("Content-Type", "image/png")
  captcha.WriteImage(ctx.ResponseWriter, captchaId, 150, 50)
}

func ServeJson(ctx *iris.Context, v interface{}) error {
  code :=  ctx.StatusCode()
  if code == 0 {
    code = iris.StatusOK;
  }
  return ctx.RenderWithStatus(code, "application/json", v)
}

