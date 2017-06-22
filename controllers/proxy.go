package controllers

import (
  "github.com/kataras/iris"
  "github.com/kataras/iris/context"
  "github.com/alsmile/goApiGateway/services/sites"
  "log"
  "gopkg.in/mgo.v2/bson"
  "net/http"
  "io/ioutil"
  "github.com/alsmile/goApiGateway/models"
  "time"
)

func ProxyDo(ctx context.Context) {
  ret := make(map[string]interface{})

  host := ctx.Host()
  method := string(ctx.Method())
  url := ctx.Params().Get("url")

  // 查找api级别代理
  siteApi, err := sites.GetApiByUrl(host, method, url)
  if err == nil {
    if siteApi.IsMock {
      ctx.WriteWithExpiration(http.StatusOK, []byte(siteApi.ResponseParamsText), siteApi.DataType, time.Now())
    } else {
      proxy(ctx, method, siteApi.Site.DstUrl+url, siteApi.DataType)
    }
    return
  }

  // 查找site级别代理
  site, err := sites.GetSiteByDomain(host)

  if err == nil {
    proxy(ctx, method, site.DstUrl+url, "")

    siteApi = &models.SiteApi{}

    // 添加到自动发现
    siteApi.AutoReg = true
    siteApi.Site.Id = site.Id
    siteApi.Method = method
    siteApi.Url = url
    siteApi.Visited = 1
    siteApi.StatusCode = ctx.GetStatusCode()
    sites.SaveApi(siteApi)

    return
  }

  // log
  ret["method"] = method
  ret["url"] = url
  ret["error"] = "Not found."
  ctx.StatusCode(iris.StatusNotFound)
  ctx.JSON(ret)
}

func proxy(ctx context.Context, method, dstUrl, dataType string) (err error) {
  client := &http.Client{}
  clientReq, err := http.NewRequest(method, dstUrl, ctx.Request().Body)
  if err != nil {
    log.Printf("[error]servers.controllers.proxy.proxy: method=%v, url=%v, proxyError=%v [[[[[[[[[[in NewRequest]]]]]]]]]]\r\n",
      method, dstUrl, err)

    ctx.StatusCode(iris.StatusNotFound)
    ctx.JSON(bson.M{"error": err.Error()})

    return
  }

  clientReq.Header = ctx.Request().Header
  clientResp, err := client.Do(clientReq)

  if err != nil {
    ctx.StatusCode(iris.StatusNotFound)
    ctx.JSON(bson.M{"error": err.Error()})
    return
  }
  ctx.StatusCode(clientResp.StatusCode)

  defer clientResp.Body.Close()
  body, err := ioutil.ReadAll(clientResp.Body)
  if err != nil {
    log.Printf("[error]servers.controllers.proxy.proxy: method=%v, url=%v, proxyError=%v, body=%v\r\n",
      method, dstUrl, err, body)
    ctx.JSON(bson.M{"error": err.Error()})

    return
  }

  ctx.WriteWithExpiration(http.StatusOK, body, dataType, time.Now())

  err = nil
  return
}

func ProxyTest(ctx context.Context) {
  method := string(ctx.Method())
  host := ctx.URLParam("host")
  url := ctx.URLParam("url")
  dataType := ctx.URLParam("dataType")

  siteApi, err := sites.GetApiByDstUrl(host, method, url)
  if err == nil && siteApi.IsMock{
    ctx.WriteWithExpiration(http.StatusOK, []byte(siteApi.ResponseParamsText), siteApi.DataType, time.Now())
    return
  }

  proxy(ctx, method, host + url, dataType)
}
