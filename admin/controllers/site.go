package controllers

import (
  "gopkg.in/kataras/iris.v6"
  "github.com/alsmile/goMicroServer/models"
  "github.com/alsmile/goMicroServer/services/user"
  "github.com/alsmile/goMicroServer/services/sites"
  "github.com/alsmile/goMicroServer/services"
  "gopkg.in/mgo.v2/bson"
)

func SiteList(ctx *iris.Context)  {
  ret := make(map[string]interface{})
  defer ServeJson(ctx, ret)

  pageIndex, err := ctx.URLParamInt(services.PageIndex)
  if err != nil || pageIndex < 1 {
    ret["error"] = services.ErrorParamPage
    return
  }
  pageCount, err := ctx.URLParamInt(services.PageCount)
  if err != nil || pageCount < 1 {
    ret["error"] = services.ErrorParamPage
    return
  }


  u := models.User{}
  user.ValidToken(ctx, &u)

  list, err := sites.List(u.Id, pageIndex, pageCount)

  if err != nil {
    ret["error"] = err.Error()
  }

  ret["list"] = list
}

func SiteGet(ctx *iris.Context)  {
  ret := make(map[string]interface{})

  id := ctx.URLParam("id")
  if id == "" {
    ret["error"] = services.ErrorParam
    ServeJson(ctx, ret)
    return
  }

  u := models.User{}
  user.ValidToken(ctx, &u)

  // 校验编辑权限
  // ...
  // End 校验编辑权限

  site := &models.Site{Id:bson.ObjectIdHex(id)}
  err := sites.Get(site)

  if err != nil {
    ret["error"] = err.Error()
    ServeJson(ctx, ret)
    return
  }

  ServeJson(ctx, site)
}

func SiteSave(ctx *iris.Context)  {
  ret := make(map[string]interface{})
  defer ServeJson(ctx, ret)

  site := &models.Site{}
  err := ctx.ReadJSON(site)
  if err != nil {
    ret["error"] = services.ErrorParam
    ret["errorConsole"] = err.Error()
  }

  u := models.User{}
  user.ValidToken(ctx, &u)

  // 校验编辑权限
  // ...
  // End 校验编辑权限

  err = sites.Save(site)

  if err != nil {
    ret["error"] = err.Error()
  }
}
