package controllers

import (
  "gopkg.in/kataras/iris.v6"
  "github.com/alsmile/goMicroServer/admin/models"
  modelsSite  "github.com/alsmile/goMicroServer/sites/models"
  "github.com/alsmile/goMicroServer/admin/services/user"
  "fmt"
  "github.com/alsmile/goMicroServer/sites/services"
)

func SiteSave(ctx *iris.Context)  {
  ret := make(map[string]interface{})
  defer ServeJson(ctx, ret)

  dataSite := &modelsSite.Site{}
  ctx.ReadJSON(dataSite)

  fmt.Printf("site=%v\r\n", dataSite)

  u := models.User{}
  user.ValidToken(ctx, &u)

  dataSite.UserId = u.Id.Hex()

  err := services.Save(dataSite)

  if err != nil {
    ret["error"] = err.Error()
  }
}
