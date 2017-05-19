package controllers
import (
  "gopkg.in/kataras/iris.v6"
  "github.com/garyburd/redigo/redis"
  "github.com/dchest/captcha"
  "github.com/alsmile/goApiGateway/session"
  myCaptcha "github.com/alsmile/goApiGateway/services/captcha"
)


func Index(ctx *iris.Context) {
  ctx.SetStatusCode(iris.StatusOK)
  ctx.ServeFile("./admin/web/dist/index.html", true)
}

func Browser(ctx *iris.Context) {
  ctx.ServeFile("./admin/web/dist/browser.html", true)
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

