package controllers

import (
  "github.com/kataras/iris"
  "github.com/kataras/iris/context"
  "github.com/alsmile/goApiGateway/utils"
  "github.com/alsmile/goApiGateway/models"
  "github.com/alsmile/goApiGateway/services/user"
  "github.com/alsmile/goApiGateway/session"
  "github.com/alsmile/goApiGateway/services/captcha"
  "github.com/alsmile/goApiGateway/services"
)

func GetSignConfig(ctx context.Context) {
  ctx.JSON(utils.GlobalConfig.User)
}

func Login(ctx context.Context) {
  ret := make(map[string]interface{})
  defer ctx.JSON(ret)

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
  ret["email"] = u.Profile.Email
  ret["username"] = u.Profile.Username
  if rememberMe {
    ret["token"] = user.GetToken(u, services.TokenValidRemember)
  } else {
    ret["token"] = user.GetToken(u, services.TokenValidHours)
  }
}


func SignUp(ctx context.Context) {
  ret := make(map[string]interface{})
  defer ctx.JSON(ret)

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

func SignActive(ctx context.Context) {
  ret := make(map[string]interface{})
  defer ctx.JSON(ret)

  u := &models.User{}
  ctx.ReadJSON(u)

  err := user.Active(u)
  if err != nil {
    ret["error"] = err.Error()
  }

  ret["id"] = u.Id
  ret["email"] = u.Profile.Email
  ret["username"] = u.Profile.Username
  ret["token"] = user.GetToken(u, services.TokenValidHours)
}

func ForgetPassword(ctx context.Context) {
  ret := make(map[string]interface{})
  defer ctx.JSON(ret)

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

func NewPassword(ctx context.Context) {
  ret := make(map[string]interface{})
  defer ctx.JSON(ret)

  u := &models.User{}
  ctx.ReadJSON(u)

  err := user.NewPassword(u)
  if err != nil {
    ret["error"] = err.Error()
    return
  }

  ret["id"] = u.Id
  ret["email"] = u.Profile.Email
  ret["username"] = u.Profile.Username
  ret["token"] = user.GetToken(u, services.TokenValidHours)
}

func UserProfile(ctx context.Context) {
  ret := make(map[string]interface{})
  defer ctx.JSON(ret)

  u := models.User{}
  user.ValidToken(ctx, &u)
  if u.Id == "" {
    ctx.StatusCode(iris.StatusUnauthorized)
    ret["error"] = services.ErrorNeedSign
    return
  }

  err := user.GetUserById(&u)
  if err != nil {
    ctx.StatusCode(iris.StatusUnauthorized)
    ret["error"] = services.ErrorUserNoExists
    return
  }

  ret["id"] = u.Id
  ret["email"] = u.Profile.Email
  ret["username"] = u.Profile.Username
}

func Auth(ctx context.Context) {
  u := models.User{}
  uid := user.ValidToken(ctx, &u)
  if uid == "" {
    ctx.StatusCode(iris.StatusUnauthorized)
    ret := make(map[string]interface{})
    ret["error"] = services.ErrorNeedSign
    ctx.JSON(ret)
    return
  }

  ctx.Next()
}
