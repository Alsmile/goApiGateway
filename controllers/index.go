package controllers

import (
	"strings"

	myCaptcha "github.com/alsmile/goApiGateway/services/captcha"
	"github.com/alsmile/goApiGateway/session"
	"github.com/dchest/captcha"
	"github.com/garyburd/redigo/redis"
	"github.com/kataras/iris"
)

// Index 首页静态文件
func Index(ctx iris.Context) {
	ctx.StatusCode(iris.StatusOK)
	ctx.ServeFile("./admin/web/dist/index.html", true)
}

// Browser 浏览器不兼容静态文件
func Browser(ctx iris.Context) {
	ctx.ServeFile("./admin/web/dist/browser.html", true)
}

// NotFound 404
func NotFound(ctx iris.Context) {
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

// Assets 静态文件
func Assets(ctx iris.Context) {
	path := ctx.Params().Get("path")
	ctx.StatusCode(iris.StatusOK)
	ctx.ServeFile("./admin/web/dist/assets/"+path, true)
}

// Captcha 验证码
func Captcha(ctx iris.Context) {
	captchaID, _ := redis.String(session.GetSession(ctx, myCaptcha.CaptchaSessionName))

	// Delete the old.
	if captchaID != "" {
		captcha.VerifyString(captchaID, "")
	}

	captchaID = captcha.New()
	session.SetSession(ctx, myCaptcha.CaptchaSessionName, captchaID)
	ctx.Header("Content-Type", "image/png")
	captcha.WriteImage(ctx.ResponseWriter(), captchaID, 150, 50)
}

// Cors 跨域设置
func Cors(ctx iris.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, token, X-Requested-With")
	ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")

	ctx.StatusCode(iris.StatusOK)
	ctx.Next()
}
