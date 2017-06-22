package controllers

import (
  "github.com/kataras/iris"
  "github.com/kataras/iris/context"
  "github.com/alsmile/goApiGateway/models"
  "github.com/alsmile/goApiGateway/services/user"
  "github.com/alsmile/goApiGateway/services/sites"
  "github.com/alsmile/goApiGateway/services"
  "gopkg.in/mgo.v2/bson"
)

func SiteList(ctx context.Context)  {
  ret := make(map[string]interface{})
  defer ctx.JSON(ret)

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

func SiteGet(ctx context.Context)  {
  ret := make(map[string]interface{})

  id := ctx.URLParam("id")
  if id == "" {
    ret["error"] = services.ErrorParam
    ctx.JSON(ret)
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
    ctx.JSON(ret)
    return
  }

  ctx.JSON(site)
}

func SiteSave(ctx context.Context)  {
  ret := make(map[string]interface{})
  defer ctx.JSON(ret)

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


func SiteApiSave(ctx context.Context)  {
  ret := make(map[string]interface{})
  defer ctx.JSON(ret)

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
      ctx.StatusCode(iris.StatusUnauthorized)
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

func SiteApiGet(ctx context.Context)  {
  ret := make(map[string]interface{})

  id := ctx.URLParam("id")
  if id == "" {
    ret["error"] = services.ErrorParam
    ctx.JSON(ret)
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
    ctx.JSON(ret)
    return
  }

  ctx.JSON(siteApi)
}

func SiteApiDel(ctx context.Context)  {
  ret := make(map[string]interface{})

  id := ctx.URLParam("id")
  if id == "" {
    ret["error"] = services.ErrorParam
    ctx.JSON(ret)
    return
  }

  u := models.User{}
  user.ValidToken(ctx, &u)

  // 校验编辑权限
  // ...
  // End 校验编辑权限

  err := sites.DelApi(id)
  if err != nil {
    ret["error"] = err.Error()
    ctx.JSON(ret)
    return
  }

  ctx.JSON(true)
}

func SiteApiList(ctx context.Context)  {
  ret := make(map[string]interface{})
  defer ctx.JSON(ret)

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

  auto := ctx.URLParam("auto")
  fieldType, _ := ctx.URLParamInt("field")

  u := models.User{}
  user.ValidToken(ctx, &u)

  siteId := bson.ObjectIdHex(ctx.URLParam("siteId"))
  list, err := sites.ApiList(siteId, auto, fieldType, pageIndex, pageCount)

  if err != nil {
    ret["error"] = err.Error()
  }

  ret["list"] = list
}

func SiteApiListByDomains(ctx context.Context)  {
  ret := make(map[string]interface{})
  defer ctx.JSON(ret)

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

  auto := ctx.URLParam("auto")
  fieldType, _ := ctx.URLParamInt("field")

  var domains []string
  err = ctx.ReadJSON(&domains)
  if err != nil || len(domains) < 1 {
    ret["error"] = services.ErrorParam
  }

  list, err := sites.ApiListByDomains(domains, auto, fieldType, pageIndex, pageCount)

  if err != nil {
    ret["error"] = err.Error()
  }

  ret["list"] = list
}
