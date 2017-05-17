package controllers

import (
  "gopkg.in/kataras/iris.v6"
  "github.com/alsmile/goMicroServer/services/sites"
  "log"
  "encoding/json"
)

func ServeJson(ctx *iris.Context, v interface{}) error {
  code :=  ctx.StatusCode()
  if code == 0 {
    code = iris.StatusOK;
  }
  return ctx.RenderWithStatus(code, "application/json", v)
}

func ProxyDo(ctx *iris.Context) {
  method := string(ctx.Method())
  key := "/" + ctx.Param("key")
  url := ctx.Param("url")
  siteApi, err :=sites.GetApiByUrl(method, key, url)
  if err != nil {
    log.Printf("[log]servers.controolers.proxy.ProxyDo: GetApiByUrl error=%v,method=%v, key=%v, url=%v\r\n",
      err, method, key, url)
  }

  ret := make(map[string]interface{})
  log.Printf("[log]servers.controolers.proxy.ProxyDo: method=%v, key=%v, url=%v, siteApi=%v\r\n",
    method, key, url, siteApi)

  // log
  ret["method"] = method
  ret["key"] = key
  ret["url"] = url
  strApi, _ := json.Marshal(siteApi)
  ret["api"] = string(strApi)
  //ret["error"] = err.Error()
  ServeJson(ctx, ret)
}
