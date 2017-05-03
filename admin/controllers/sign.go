package controllers

import (
  "gopkg.in/kataras/iris.v6"
  "github.com/alsmile/goMicroServer/utils"
  "github.com/alsmile/goMicroServer/admin/models"
  "github.com/alsmile/goMicroServer/admin/services/user"
  "github.com/alsmile/goMicroServer/session"
  "github.com/alsmile/goMicroServer/services/captcha"
  "github.com/alsmile/goMicroServer/services"
)

func GetSignConfig(ctx *iris.Context) {
  ServeJson(ctx, utils.GlobalConfig.User)
}

func Login(ctx *iris.Context) {
  ret := make(map[string]interface{})
  defer ServeJson(ctx, ret)

  u := &models.User{}
  ctx.ReadJSON(u)
  rememberMe := u.RememberMe

  sid := session.GetSessionId(ctx)
  if captcha.IsNeedSignCaptcha(sid) {
    if captcha.VerifyImage(ctx, u.Captcha) == false {
      ret["error"] = services.ErrorCaptchaCode
      ret["errorTip"] = "captcha"
      return
    }
  }

  err := user.GetUserByPassword(u)
  if err != nil {
    ret["error"] = err.Error()
    captcha.SignError(sid)
    if captcha.IsNeedSignCaptcha(sid) {
      ret["errorTip"] = "captcha"
    }
    return
  }

  ret["id"] = u.Id
  ret["email"] = u.Email
  ret["username"] = u.Username
  if rememberMe {
    ret["token"] = user.GetToken(u, services.TokenValidRemember)
  } else {
    ret["token"] = user.GetToken(u, services.TokenValidHours)
  }
}


func SignUp(ctx *iris.Context) {
  ret := make(map[string]interface{})
  defer ServeJson(ctx, ret)

  u := &models.User{}
  ctx.ReadJSON(u)

  if captcha.VerifyImage(ctx, u.Captcha) == false {
    ret["error"] = services.ErrorCaptchaCode
    return
  }

  err := user.AddUser(u)
  if err != nil {
    ret["error"] = err.Error()
  }
}

func SignActive(ctx *iris.Context) {
  ret := make(map[string]interface{})
  defer ServeJson(ctx, ret)

  u := &models.User{}
  ctx.ReadJSON(u)

  err := user.Active(u)
  if err != nil {
    ret["error"] = err.Error()
  }

  ret["id"] = u.Id
  ret["email"] = u.Email
  ret["username"] = u.Username
  ret["token"] = user.GetToken(u, services.TokenValidHours)
}

func ForgetPassword(ctx *iris.Context) {
  ret := make(map[string]interface{})
  defer ServeJson(ctx, ret)

  u := &models.User{}
  ctx.ReadJSON(u)

  if captcha.VerifyImage(ctx, u.Captcha) == false {
    ret["error"] = services.ErrorCaptchaCode
    return
  }

  err := user.ForgetPassword(u)
  if err != nil {
    ret["error"] = err.Error()
  }
}

func NewPassword(ctx *iris.Context) {
  ret := make(map[string]interface{})
  defer ServeJson(ctx, ret)

  u := &models.User{}
  ctx.ReadJSON(u)

  err := user.NewPassword(u)
  if err != nil {
    ret["error"] = err.Error()
    return
  }

  ret["id"] = u.Id
  ret["email"] = u.Email
  ret["username"] = u.Username
  ret["token"] = user.GetToken(u, services.TokenValidHours)
}

func UserInfo(ctx *iris.Context) {
  ret := make(map[string]interface{})
  defer ServeJson(ctx, ret)

  u := models.User{}
  user.ValidToken(ctx, &u)
  if u.Id == "" {
    ret["error"] = services.ErrorNeedSign
    return
  }

  err := user.GetUserById(&u)
  if err != nil {
    ret["error"] = services.ErrorUserNoExists
    return
  }

  ret["id"] = u.Id
  ret["email"] = u.Email
  ret["username"] = u.Username
}
