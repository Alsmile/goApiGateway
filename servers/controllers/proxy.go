package controllers

import (
  "gopkg.in/kataras/iris.v6"
  "github.com/alsmile/goApiGateway/services/sites"
  "log"
  "encoding/json"
  "gopkg.in/mgo.v2/bson"
  "net/http"
  "io/ioutil"
  "github.com/alsmile/goApiGateway/utils"
)

func ServeJson(ctx *iris.Context, v interface{}) error {
  code := ctx.StatusCode()
  if code == 0 {
    code = iris.StatusOK;
  }
  return ctx.RenderWithStatus(code, "application/json", v)
}

func ProxyDo(ctx *iris.Context) {
  ret := make(map[string]interface{})

  subdomain := ctx.Subdomain()
  if ctx.VirtualHostname() == utils.GlobalConfig.Domain.Domain {
    subdomain = ""
  }
  method := string(ctx.Method())
  key := "/" + ctx.Param("key")
  url := ctx.Param("url")

  // 查找api级别代理
  siteApi, err := sites.GetApiByUrl(subdomain, method, key, url)
  if err == nil {
    if siteApi.IsMock {
      if siteApi.DataType == "application/json" ||
        siteApi.DataType == "multipart/form-data" ||
        siteApi.DataType == "application/x-www-form-urlencoded" {
        var data bson.M
        json.Unmarshal([]byte(siteApi.ResponseParamsText), &data)
        ctx.RenderWithStatus(iris.StatusOK, siteApi.DataType, data)
      } else {
        ctx.RenderWithStatus(iris.StatusOK, siteApi.DataType, siteApi.ResponseParamsText)
      }
    } else {
      proxy(ctx, method, siteApi.Site.ProxyValue+url, siteApi.DataType)
    }
    return
  }

  // 查找api级别代理
  site, err := sites.GetSiteByProxyKey(subdomain, key)
  if err == nil {
    proxy(ctx, method, site.ProxyValue+url, "")
    return
  }

  // log
  ret["method"] = method
  ret["key"] = key
  ret["url"] = url
  ret["error"] = "Not found."
  ctx.SetStatusCode(iris.StatusNotFound)
  ServeJson(ctx, ret)
}

func proxy(ctx *iris.Context, method, dstUrl, dataType string) (err error) {
  client := &http.Client{}
  clientReq, err := http.NewRequest(method, dstUrl, ctx.Request.Body)
  if err != nil {
    log.Printf("[error]servers.controllers.proxy.proxy: method=%v, url=%v, proxyError=%v [[[[[[[[[[in NewRequest]]]]]]]]]]\r\n",
      method, dstUrl, err)

    ctx.SetStatusCode(iris.StatusNotFound)
    ServeJson(ctx, bson.M{"error": err.Error()})

    return
  }

  clientReq.Header = ctx.Request.Header
  clientResp, err := client.Do(clientReq)

  if err != nil {
    ctx.SetStatusCode(iris.StatusNotFound)
    ServeJson(ctx, bson.M{"error": err.Error()})
    return
  }
  ctx.SetStatusCode(clientResp.StatusCode)

  defer clientResp.Body.Close()
  body, err := ioutil.ReadAll(clientResp.Body)
  if err != nil {
    log.Printf("[error]servers.controllers.proxy.proxy: method=%v, url=%v, proxyError=%v, body=%v\r\n",
      method, dstUrl, err, body)
    ServeJson(ctx, bson.M{"error": err.Error()})

    return
  }

  if dataType == "application/json" ||
    dataType == "multipart/form-data" ||
    dataType == "application/x-www-form-urlencoded" {
    var data bson.M
    json.Unmarshal(body, &data)
    ctx.Render(dataType, data)
  } else if dataType == "" {
    var data bson.M
    err = json.Unmarshal(body, &data)
    if err == nil {
      ctx.Render(dataType, data)
    } else {
      ctx.Render(dataType, string(body))
    }
  } else {
    ctx.Render(dataType, string(body))
  }

  err = nil
  return
}

func ProxyTest(ctx *iris.Context) {
  method := string(ctx.Method())
  url := ctx.URLParam("url")
  dataType := ctx.URLParam("dataType")

  proxy(ctx, method, url, dataType)
}
