package controllers

import (
  "gopkg.in/kataras/iris.v6"
  "github.com/alsmile/goMicroServer/utils"
  "github.com/alsmile/goMicroServer/admin/models"
  "github.com/alsmile/goMicroServer/admin/services/user"
)

func GetSignConfig(ctx *iris.Context) {
  ServeJson(ctx, utils.GlobalConfig.User)
}

func Login(ctx *iris.Context) {
  ret := make(map[string]interface{})
  defer ServeJson(ctx, ret)

  u := &models.User{}
  ctx.ReadJSON(u)

  err := user.GetUserByPassword(u)
  if err != nil {
    ret["error"] = err.Error()
    return
  }

  ret["id"] = u.Id
  ret["email"] = u.Email
  ret["username"] = u.Username
  ret["token"] = user.GetToken(u, 10)
}
