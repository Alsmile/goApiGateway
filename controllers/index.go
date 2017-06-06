package controllers
import (
  "github.com/kataras/iris"
  "github.com/kataras/iris/context"
  "github.com/garyburd/redigo/redis"
  "github.com/dchest/captcha"
  "github.com/alsmile/goApiGateway/session"
  myCaptcha "github.com/alsmile/goApiGateway/services/captcha"
  "strings"
)


func Index(ctx context.Context) {
  ctx.StatusCode(iris.StatusOK)
  ctx.ServeFile("./admin/web/dist/index.html", true)
}

func Browser(ctx context.Context) {
  ctx.ServeFile("./admin/web/dist/browser.html", true)
}

func NotFound(ctx context.Context) {
  if strings.HasPrefix(ctx.Path(), "/api/") {
    ret := make(map[string]interface{})
    ret["error"] = "请求错误（Not found）：" + ctx.Path()
    ctx.StatusCode(iris.StatusNotFound)
    ctx.JSON(ret)
  } else {
    ctx.StatusCode(iris.StatusFound)
    Index(ctx)
  }
}

func Captcha(ctx context.Context) {
  captchaId, _ := redis.String(session.GetSession(ctx, myCaptcha.CaptchaSessionName))

  // Delete the old.
  if captchaId != "" {
    captcha.VerifyString(captchaId, "")
  }

  captchaId = captcha.New()
  session.SetSession(ctx, myCaptcha.CaptchaSessionName, captchaId)
  ctx.Header("Content-Type", "image/png")
  captcha.WriteImage(ctx.ResponseWriter(), captchaId, 150, 50)
}


