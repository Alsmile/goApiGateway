package controllers

import (
  "time"
  "encoding/json"
  "gopkg.in/mgo.v2/bson"
  "github.com/kataras/iris/context"
  "github.com/alsmile/goApiGateway/services/sites"
  "github.com/alsmile/goApiGateway/models"
  "github.com/alsmile/goApiGateway/services"
)

func ApisSet(ctx context.Context) {
  ret := make(map[string]interface{})
  defer ctx.JSON(ret)

  sdkSite := &models.SdkSite{}
  err := ctx.ReadJSON(sdkSite)
  if err != nil {
    ret["error"] = services.ErrorParam
    return
  }

  if sdkSite.Name == "" {
    sdkSite.Name = sdkSite.ApiDomain
  }

  site := &models.Site{
    Name: sdkSite.Name,
    Gzip: sdkSite.Gzip,
    Https: sdkSite.Https,
    ApiDomain: sdkSite.ApiDomain,
    Group: sdkSite.Group,
    DstUrl: sdkSite.DstUrl,
  }

  // 先删除
  sites.DelSite(site)

  // 再保存site信息
  err = sites.Save(site)
  if err != nil {
    ret["error"] = err.Error()
    return
  }

  var apis []models.SiteApi
  json.Unmarshal([]byte(sdkSite.Apis), &apis)

  if len(apis) < 1 {
    return
  }

  docs := make([]interface{}, len(apis))
  for i:=0; i < len(apis); i++ {
    apis[i].Site.Id = site.Id
    apis[i].Site.ApiDomain = site.ApiDomain
    apis[i].Site.DstUrl = site.DstUrl
    apis[i].Site.Group = site.Group
    apis[i].Site.Https = site.Https
    apis[i].Site.Gzip = site.Gzip
    apis[i].Id = bson.NewObjectId()
    apis[i].UpdatedAt = time.Now().UTC()
    apis[i].CreatedAt = apis[i].UpdatedAt

    docs[i] = apis[i]
  }
  sites.SaveApis(docs)
}

func ApisDelete(ctx context.Context) {
  ret := make(map[string]interface{})
  defer ctx.JSON(ret)

  sdkSite := &models.SdkSite{}
  err := ctx.ReadJSON(sdkSite)
  if err != nil {
    ret["error"] = services.ErrorParam
    return
  }

  site := &models.Site{
    Name: sdkSite.Name,
    Gzip: sdkSite.Gzip,
    Https: sdkSite.Https,
    ApiDomain: sdkSite.ApiDomain,
    Group: sdkSite.Group,
    DstUrl: sdkSite.DstUrl,
  }
  // 删除site信息
  err = sites.DelSite(site)
  if err != nil {
    ret["error"] = err.Error()
    return
  }
}
