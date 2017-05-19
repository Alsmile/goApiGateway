package controllers

import (
  "gopkg.in/kataras/iris.v6"
  "github.com/alsmile/goApiGateway/models"
  "github.com/alsmile/goApiGateway/services/user"
  "github.com/alsmile/goApiGateway/services/sites"
  "github.com/alsmile/goApiGateway/services"
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


func SiteApiSave(ctx *iris.Context)  {
  ret := make(map[string]interface{})
  defer ServeJson(ctx, ret)

  siteApi := &models.SiteApi{}
  err := ctx.ReadJSON(siteApi)
  if err != nil {
    ret["error"] = services.ErrorParam
    ret["errorConsole"] = err.Error()
  }

  u := models.User{}
  user.ValidToken(ctx, &u)

  //  site id不存在，表示自动根据api信息保存site
  if siteApi.Id == "" && siteApi.Site.Id == "" {
    err = user.GetUserById(&u)
    if err != nil {
      ctx.SetStatusCode(iris.StatusUnauthorized)
      ret["error"] = services.ErrorUserNoExists
      return
    }

    siteApi.Owner.Id = u.Id
    siteApi.Owner.Email = u.Profile.Email
    siteApi.Owner.Phone = u.Profile.Phone
    siteApi.Owner.Username = u.Profile.Username
    siteApi.Editor = siteApi.Owner
  }

  err = sites.SaveApi(siteApi)
  if err != nil {
    ret["error"] = err.Error()
  }
  ret["id"] = siteApi.Id
}

func SiteApiGet(ctx *iris.Context)  {
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

  siteApi := &models.SiteApi{Id:bson.ObjectIdHex(id)}
  err := sites.GetApi(siteApi)

  if err != nil {
    ret["error"] = err.Error()
    ServeJson(ctx, ret)
    return
  }

  ServeJson(ctx, siteApi)
}

func SiteApiList(ctx *iris.Context)  {
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

  siteId := bson.ObjectIdHex(ctx.URLParam("siteId"))
  list, err := sites.ApiList(siteId, pageIndex, pageCount)

  if err != nil {
    ret["error"] = err.Error()
  }

  ret["list"] = list
}
